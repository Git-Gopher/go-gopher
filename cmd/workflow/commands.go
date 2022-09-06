package main

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/model"
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
)

var _ Commands = &Cmds{}

type Commands interface {
	SingleUrlCommand(cCtx *cli.Context, flags *Flags) error
	SingleLocalCommand(cCtx *cli.Context, flags *Flags) error
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

func (c *Cmds) runRules(repo *git.Repository, githubURL string) error {
	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return fmt.Errorf("failed to get owner and repo name: %w", err)
	}

	// Create enrichedModel.
	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("failed to create enriched model: %w", err)
	}

	scoresMap := workflow.Detect(rule.RuleCtx{
		Model:          enrichedModel,
		LoginWhiteList: []string{}, // TODO: add whitelist
	})

	log.Infof("Finished running rules for %s/%s", repoOwner, repoName)

	for name, scores := range scoresMap {
		log.Infof("Rule: %s", name)
		log.Infof("	GitHubFlow: %v", scores.GitHubFlow().Value())
		log.Infof("	GitFlow: %v", scores.GitFlow().Value())
		log.Infof("	GitLabFlow: %v", scores.GitLabFlow().Value())
		log.Infof("	OneFlow: %v", scores.OneFlow().Value())
	}

	return nil
}
