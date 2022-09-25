package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	errGitHubURL = fmt.Errorf("missing GitHub URL")
	errLocalDir  = fmt.Errorf("missing Local Directory")
	errBatchJson = fmt.Errorf("missing repository json file")
)

type logs struct {
	Url string `json:"url"`
	// If the repository has been skipped from running due to timeout.
	// Skipped          bool                    `json"skipped"`
	Scores           map[string]*rule.Scores `json:"scores"`
	DetectedWorkflow []string                `json:"detectedWorkflow"`
}

var _ Commands = &Cmds{}

type Commands interface {
	SingleUrlCommand(cCtx *cli.Context, flags *Flags) error
	SingleLocalCommand(cCtx *cli.Context, flags *Flags) error
	BatchUrlCommand(cCtx *cli.Context, flags *Flags) error
}

type Cmds struct{}

func (c *Cmds) SingleUrlCommand(cCtx *cli.Context, flags *Flags) error {
	githubURL := cCtx.Args().Get(0)
	if githubURL == "" {
		return errGitHubURL
	}

	var auth *githttp.BasicAuth
	if flags.GithubToken != "" {
		auth = &githttp.BasicAuth{
			Username: "non-empty",
			Password: flags.GithubToken,
		}
	}

	// Clone repository into memory.
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  githubURL,
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	if err = c.runRules(repo, githubURL); err != nil {
		return err
	}

	return nil
}

func (c *Cmds) SingleLocalCommand(cCtx *cli.Context, flags *Flags) error {
	directory := cCtx.Args().Get(0)
	if directory == "" {
		return errLocalDir
	}

	return c.runLocalRepository(directory)
}

func (c *Cmds) runLocalRepository(directory string) error {
	// Open repository locally.
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	githubURL, err := utils.Url(repo)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}

	if err = c.runRules(repo, githubURL); err != nil {
		return err
	}

	return nil
}

// nolint: gocognit
func (c *Cmds) BatchUrlCommand(cCtx *cli.Context, flags *Flags) error {
	repositoryJsonPath := cCtx.Args().Get(0)
	if repositoryJsonPath == "" {
		return errBatchJson
	}

	payload, err := os.ReadFile(filepath.Clean(repositoryJsonPath))
	if err != nil {
		return fmt.Errorf("unable to read json file payload: %w", err)
	}

	var repositories []remote.Repository
	if err = json.Unmarshal(payload, &repositories); err != nil {
		return fmt.Errorf("unable to marshal payload to repositories: %w", err)
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	repoUrlChan := make(chan string, len(repositories))

	wg.Add(1)
	go func() {
		for _, repo := range repositories {
			wg.Add(1)
			repoUrlChan <- repo.Url
		}
		wg.Done()
	}()

	log.Infof("Running batch workflow detection on %d threads...", runtime.NumCPU())
	for i := 0; i < runtime.NumCPU()-1; i++ {
		go func() {
			for {
				select {
				case url := <-repoUrlChan:
					var auth *githttp.BasicAuth

					if flags.GithubToken != "" {
						auth = &githttp.BasicAuth{
							Username: "non-empty",
							Password: flags.GithubToken,
						}
					} else {
						log.Errorf("no github token passed in as flag (required for upgraded api limits): %v", err)
						wg.Done()

						return
					}

					log.Infof("Cloning repository %s to memory...", url)

					start := time.Now()
					// nolint: contextcheck
					repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
						URL:  url,
						Auth: auth,
					})
					if err != nil {
						log.Errorf("failed to clone repository: %v", err)
						wg.Done()

						return
					}

					log.Infof("Finished repository %s to memory (%s)...", url, time.Since(start))
					if err = c.runRules(repo, url); err != nil {
						log.Errorf("failed to run rules: %v", err)
						wg.Done()

						return
					}

					wg.Done()

				case <-time.After(time.Duration(flags.Timeout) * time.Second):
					log.Error("Repository timed out, continuing...")
					wg.Done()

					wg.Done()

				case <-ctx.Done():
					log.Error("context finished")

					return
				}

				log.Infof("Repositories left to process: %d", len(repoUrlChan))
				if len(repoUrlChan) == 0 {
					break
				}
			}
		}()
	}

	wg.Wait()
	cancel()

	return nil
}

func (c *Cmds) runRules(repo *git.Repository, githubURL string) error {
	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return fmt.Errorf("failed to get owner and repo name: %w", err)
	}

	// Create enrichedModel.
	log.Infof("Fetching enriched model for repository %s/%s...", repoOwner, repoName)
	start := time.Now()

	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("failed to create enriched model: %w", err)
	}
	log.Infof("Done Fetching enriched model for repository  %s/%s (%s)...", repoOwner, repoName, time.Since(start))

	log.Infof("Running rules for %s/%s", repoOwner, repoName)
	start = time.Now()

	scoresMap := workflow.Detect(rule.RuleCtx{
		Model:          enrichedModel,
		LoginWhiteList: []string{}, // TODO: add whitelist
	})

	log.Infof("Finished running rules for %s/%s (%s)", repoOwner, repoName, time.Since(start))

	accScores := make(map[rule.WorkflowType]float64) // workflow types iota
	for _, scores := range scoresMap {
		accScores[rule.GitHubFlow] += scores.GitHubFlow.Value()
		accScores[rule.GitFlow] += scores.GitFlow.Value()
		accScores[rule.GitlabFlow] += scores.GitLabFlow.Value()
		accScores[rule.OneFlow] += scores.OneFlow.Value()
		accScores[rule.TrunkBased] += scores.TrunkBased.Value()
	}

	// Detect workflow type.
	// Slice as might be multiple equal scores.
	var detectedWorkflow []string
	{
		max := float64(0)
		for k, v := range accScores {
			if v > max {
				max = v
				detectedWorkflow = []string{k.String()}
			} else if v == max {
				detectedWorkflow = append(detectedWorkflow, k.String())
			}
		}
	}

	if err = writeLog(githubURL, scoresMap, detectedWorkflow, repoOwner, repoName); err != nil {
		return err
	}

	return nil
}

// Write an individual log file for the repository to disk.
func writeLog(githubURL string,
	scoresMap map[string]*rule.Scores,
	detectedWorkflow []string,
	repoOwner string,
	repoName string,
) error {
	logFilePath := fmt.Sprintf("output/workflow-output-%s-%s.json", repoOwner, repoName)

	log.Infof("Writing log for %s/%s to %s...", repoOwner, repoName, logFilePath)
	logging := logs{
		Url:              githubURL,
		DetectedWorkflow: detectedWorkflow,
		Scores:           scoresMap,
		// Skipped:          false,
	}

	payload, err := json.MarshalIndent(logging, "", " ")
	if err != nil {
		return fmt.Errorf("could not marshal log payload: %w", err)
	}

	if err = os.WriteFile(logFilePath, payload, 0o600); err != nil {
		return fmt.Errorf("failed to write log file: %w", err)
	}

	log.Infof("Output workflow logs to for %s/%s to %s", repoOwner, repoName, logFilePath)

	return nil
}
