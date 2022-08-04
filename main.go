package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/config"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v45/github"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

//nolint
func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = version.BuildVersion()
	app.Commands = []*cli.Command{
		{
			Name:    "action",
			Aliases: []string{"a"},
			Usage:   "detect a workflow for current root",
			Action: func(ctx *cli.Context) error {
				utils.Environment(".env")
				workspace := utils.EnvGithubWorkspace()

				// repository := os.Getenv("GITHUB_REPOSITORY")
				// sha := os.Getenv("GITHUB_SHA") // commit sha triggered
				// ref := os.Getenv("GITHUB_REF") // branch ref triggered

				// Repo
				repo, err := git.PlainOpen(workspace)
				if err != nil {
					log.Fatalf("cannot read repo: %v\n", err)
				}

				gitModel, err := local.NewGitModel(repo)
				if err != nil {
					log.Fatalf("Could not create GitModel: %v\n", err)
				}

				url := os.Getenv("GITHUB_URL")
				if url == "" {
					url, err = utils.Url(repo)
					if err != nil {
						log.Fatalf("Could get url from repository: %v\n", err)
					}
					log.Printf("GITHUB_URL is not set, falling back to \"%s\"...\n", url)
				}
				owner, name, err := utils.OwnerNameFromUrl(url)
				if err != nil {
					log.Fatalf("Could not get owner and name from URL: %v\n", err)
				}

				remoteModel, err := remote.ScrapeRemoteModel(owner, name)
				if err != nil {
					log.Fatalf("Could not create RemoteModel: %v\n", err)
				}

				enrichedModel := enriched.NewEnrichedModel(*gitModel, *remoteModel)

				// Authors
				authors := enriched.PopulateAuthors(enrichedModel)

				// Cache
				current := cache.NewCache(enrichedModel)
				caches, err := cache.Read()

				//nolint
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Cache file does not exist: %v", err)
					// Write a cache for current so that next run can use it
					if err = cache.Write([]*cache.Cache{current}); err != nil {
						log.Fatalf("Could not write cache: %v\n", err)
					}
				} else if err != nil {
					log.Fatalf("Failed to load caches: %v", err)
				}

				cfg := readConfig(ctx)
				ghwf := workflow.GithubFlowWorkflow(cfg)
				violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
				if err != nil {
					log.Fatalf("Failed to analyze: %v\n", err)
				}

				workflowSummary(authors, violated, count, total, violations)

				// Set action outputs to a markdown summary.
				summary := markdownSummary(authors, violations)
				markup.Outputs("pr_summary", summary)

				err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
				if err != nil {
					log.Fatalf("Could not create csv summary: %v", err)
				}

				err = ghwf.WriteLog(*enrichedModel, cfg)
				if err != nil {
					log.Fatalf("Could not write json log: %v", err)
				}

				return nil
			},
		},
		{
			Name:  "download",
			Usage: "Download artifact logs from workflow runs for a batch of repositories",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "output",
					Aliases:     []string{"o"},
					Usage:       "output directory",
					DefaultText: "output",
					Value:       "output",
					Required:    false,
				},
			},
			Action: func(ctx *cli.Context) error {
				utils.Environment(".env")
				org := ctx.Args().Get(0)
				out := ctx.String("output")
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
				)
				tc := oauth2.NewClient(ctx.Context, ts)
				client := github.NewClient(tc)

				var logs []interface{}
				err := os.MkdirAll(out, 0o755)
				if err != nil {
					log.Fatalf("Failed to create output directory: %v", err)
				}

				log.Printf("Fetching repositories for organisation %s...\n", org)
				repos, _, err := client.Repositories.ListByOrg(ctx.Context, org, nil)
				if err != nil {
					log.Fatalf("Could not fetch orginisation repositories: %v", err)
				}

				for _, r := range repos {
					arts, _, err := client.Actions.ListArtifacts(ctx.Context, org, *r.Name, nil)
					if err != nil {
						log.Fatalf("Could not fetch artifact list: %v", err)
					}

					for _, a := range arts.Artifacts {
						log.Printf("Downloading artifacts for %s/%s...\n", org, *r.Name)
						url, _, err := client.Actions.DownloadArtifact(ctx.Context, org, *r.Name, *a.ID, true)
						if err != nil {
							log.Fatalf("could not fetch artifact url: %v", err)
						}

						pathZip := fmt.Sprintf("%s/log-%s-%s-%d.zip", out, org, *r.Name, *a.ID)
						pathJson := fmt.Sprintf("%s/log-%s-%s-%d", out, org, *r.Name, *a.ID)

						log.Printf("Downloading artifact %s...\n", pathZip)
						err = utils.DownloadFile(pathZip, url.String())
						if err != nil {
							log.Fatalf("could not download artifact: %v", err)
						}

						log.Printf("Unzipping %s...\n", pathZip)
						utils.Unzip(pathZip, pathJson)
					}
				}

				logCount := 0
				filepath.Walk(out, func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}

					if matched, err := filepath.Match("log-go-gopher*.json", filepath.Base(path)); err != nil {
						return err
					} else if matched {
						logCount++

						log.Printf("Appending %s to merged log...", path)
						file, err := ioutil.ReadFile(path)
						if err != nil {
							log.Fatalf("Could not read log file: %v", err)
						}

						var data interface{}
						json.Unmarshal(file, &data)
						logs = append(logs, data)
					}
					return nil
				})

				bytes, err := json.MarshalIndent(logs, "", " ")
				if err != nil {
					return fmt.Errorf("error marshaling merged log: %w", err)
				}

				logPath := fmt.Sprintf("%s/merged-log-%s.json", out, org)
				if err := ioutil.WriteFile(logPath, bytes, 0o600); err != nil {
					return fmt.Errorf("error writing merged log: %w", err)
				}

				log.Printf("Downloaded %d logs from %s and merged to %s", logCount, "Git-Gopher", logPath)

				return nil
			},
		},
		{
			Name:  "analyze",
			Usage: "detect a workflow for current root",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Usage:    "path to configuation file",
					Aliases:  []string{"c"},
					Required: false,
				},
				// TODO: Complete this feature
				&cli.BoolFlag{
					Name:     "logging",
					Usage:    "enable logging, output to the set file name",
					Aliases:  []string{"l"},
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "csv",
					Usage:    "csv summary of the workflow run",
					Required: false,
				},
			},

			Subcommands: []*cli.Command{
				{
					Name:  "url",
					Usage: "analyze github project url",
					Action: func(ctx *cli.Context) error {
						utils.Environment(".env")
						token := os.Getenv("GITHUB_TOKEN")
						if token == "" {
							log.Fatal("GITHUB_TOKEN is not set")
						}

						url := ctx.Args().Get(0)

						repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
							URL: url,
							Auth: &githttp.BasicAuth{
								Username: "non-empty",
								Password: token,
							},
						})
						if err != nil {
							log.Fatalf("Failed to clone repository: %v", err)
						}

						gitModel, err := local.NewGitModel(repo)
						if err != nil {
							log.Fatalf("Could not create GitModel: %v\n", err)
						}

						owner, name, err := utils.OwnerNameFromUrl(url)
						if err != nil {
							log.Fatalf("Could not get owner and name from URL: %v\n", err)
						}

						githubModel, err := remote.ScrapeRemoteModel(owner, name)
						if err != nil {
							log.Fatalf("Could not scrape GithubModel: %v\n", err)
						}
						enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

						// Authors
						authors := enriched.PopulateAuthors(enrichedModel)

						// Cache
						current := cache.NewCache(enrichedModel)
						caches, err := cache.Read()

						//nolint
						if errors.Is(err, os.ErrNotExist) {
							log.Printf("Cache file does not exist: %v", err)
							// Write a cache for current so that next run can use it
							if err = cache.Write([]*cache.Cache{current}); err != nil {
								log.Fatalf("Could not write cache: %v\n", err)
							}
						} else if err != nil {
							log.Fatalf("Failed to load caches: %v", err)
						}

						cfg := readConfig(ctx)
						ghwf := workflow.GithubFlowWorkflow(cfg)
						violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
						if err != nil {
							log.Fatalf("Failed to analyze: %v\n", err)
						}

						workflowSummary(authors, violated, count, total, violations)

						if ctx.Bool("csv") {
							err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
							if err != nil {
								log.Fatalf("Could not create csv summary: %v", err)
							}
						}

						if ctx.Bool("logging") {
							err = ghwf.WriteLog(*enrichedModel, cfg)
							if err != nil {
								log.Fatalf("Could not write json log: %v", err)
							}
						}

						return nil
					},
				},
				{
					Name:    "batch",
					Aliases: []string{"b"},
					Usage:   "batch analyze a series of local git projects",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:     "csv",
							Usage:    "csv summary of the workflow run",
							Required: false,
						},
					},
					Action: func(ctx *cli.Context) error {
						utils.Environment(".env")
						path := ctx.Args().Get(0)
						if path == "" {
							path = "./"
							log.Printf("No path provided, using current directory (\"%s\") as target path", path)

						}

						var ps []string
						filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
							if info.IsDir() && info.Name() == ".git" {
								ps = append(ps, filepath.Dir(path))
								return filepath.SkipDir
							}

							return nil
						})
						if len(ps) == 0 {
							log.Fatalf("Could not detect any git repositories within the directiory: \"%s\"", path)
						}

						cfg := readConfig(ctx)
						ghwf := workflow.GithubFlowWorkflow(cfg)

						for _, p := range ps {
							repo, err := git.PlainOpen(p)
							if err != nil {
								log.Fatalf("Failed to clone repository: %v", err)
							}

							gitModel, err := local.NewGitModel(repo)
							if err != nil {
								log.Fatalf("Could not create GitModel: %v\n", err)
							}

							url, err := utils.Url(repo)
							if err != nil {
								log.Fatalf("Could get url from repository: \"%v\", does it have any remotes?", err)
							}

							owner, name, err := utils.OwnerNameFromUrl(url)
							if err != nil {
								log.Fatalf("Could get the owner and name from URL: %v", err)
							}

							githubModel, err := remote.ScrapeRemoteModel(owner, name)
							if err != nil {
								log.Fatalf("Could not create GithubModel: %v\n", err)
							}

							enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

							// Authors
							authors := enriched.PopulateAuthors(enrichedModel)

							v, c, t, vs, err := ghwf.Analyze(enrichedModel, authors, nil, nil)
							if err != nil {
								log.Fatalf("Failed to analyze: %v\n", err)
							}

							workflowSummary(authors, v, c, t, vs)
							nameCsv := fmt.Sprintf("batch-%s.csv", filepath.Base(path))
							if ctx.Bool("csv") {
								err = ghwf.Csv(nameCsv, enrichedModel.Name, enrichedModel.URL)
								if err != nil {
									log.Fatalf("Could not create csv summary: %v", err)
								}
							}

							if ctx.Bool("logging") {
								err = ghwf.WriteLog(*enrichedModel, cfg)
								if err != nil {
									log.Fatalf("Could not write json log: %v", err)
								}
							}
						}
						return nil
					},
				},
			},
		},
	}

	log.Printf("BuildVersion: %s", version.BuildVersion())
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Print violation summary to IO, Split by severity with author association.
func workflowSummary(authors utils.Authors, v, c, t int, vs []violation.Violation) {
	var violations, suggestions []violation.Violation
	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			violations = append(violations, v)
		case violation.Suggestion:
			suggestions = append(suggestions, v)
		}
	}

	var vSb strings.Builder
	for _, v := range violations {
		vSb.WriteString(v.Display(authors))
	}
	markup.Group("Violations", vSb.String())

	var sSb strings.Builder
	for _, v := range suggestions {
		sSb.WriteString(v.Display(authors))
	}
	markup.Group("Suggestions", sSb.String())

	var aSb strings.Builder
	counts := make(map[string]int)
	for _, v := range vs {
		email := v.Email()
		login, err := authors.Find(email)
		if err != nil {
			continue
		}
		counts[*login]++
	}

	for login, count := range counts {
		aSb.WriteString(fmt.Sprintf("%s: %d\n", login, count))
	}

	aSb.WriteString(fmt.Sprintf("violated: %d\n", v))
	aSb.WriteString(fmt.Sprintf("count: %d\n", c))
	aSb.WriteString(fmt.Sprintf("total: %d\n", t))
	markup.Group("Summary", aSb.String())
}

// Fetch custom or default config. Fatal on bad custom config.
func readConfig(ctx *cli.Context) *config.Config {
	var cfg *config.Config
	var err error

	// Custom config
	if ctx.String("config") != "" {
		cfg, err = config.Read(ctx.String("config"))
		if err != nil {
			log.Fatalf("Failed to read custom config: %v", err)
		}
	} else {
		// Use default config
		cfg, err = config.Default()
		if err != nil {
			log.Fatalf("Failed to read default config: %v", err)
		}
	}

	return cfg
}

// Helper function to create a markdown summary of the violations.
func markdownSummary(authors utils.Authors, vs []violation.Violation) string {
	md := markup.CreateMarkdown("Workflow Summary")
	md.AddLine(fmt.Sprintf("Created with git-gopher version `%s`", version.BuildVersion()))

	// Separate violation types.
	var violations []violation.Violation
	var suggestions []violation.Violation

	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			violations = append(violations, v)
		case violation.Suggestion:
			suggestions = append(suggestions, v)
		default:
			log.Printf("Unknown violation severity: %v", v.Severity())
		}
	}

	headers := []string{"Violation", "Message", "Advice", "Author"}
	rows := make([][]string, len(violations))

	for i, v := range violations {
		row := make([]string, len(headers))
		name := v.Name()
		row[0] = name
		message := v.Message()
		row[1] = message

		suggestion, err := v.Suggestion()
		if err != nil {
			suggestion = ""
		}
		row[2] = suggestion

		usernamePtr, err := authors.Find(v.Email())
		if err != nil || usernamePtr == nil {
			row[3] = "@unknown"
		} else {
			row[3] = markup.Author(*usernamePtr).Markdown()
		}

		rows[i] = row
	}

	md.BeginCollapsable("Violations")
	md.Table(headers, rows)
	md.EndCollapsable()

	headers = []string{"Suggestion", "Message", "Advice", "Author"}
	rows = make([][]string, len(suggestions))

	for i, v := range suggestions {
		row := make([]string, len(headers))
		name := v.Name()
		row[0] = name
		message := v.Message()
		row[1] = message

		suggestion, err := v.Suggestion()
		if err != nil {
			suggestion = ""
		}
		row[2] = suggestion

		usernamePtr, err := authors.Find(v.Email())
		if err != nil || usernamePtr == nil {
			row[3] = "@unknown"
		} else {
			row[3] = markup.Author(*usernamePtr).Markdown()
		}

		rows[i] = row
	}

	md.BeginCollapsable("Suggestions")
	md.Table(headers, rows)
	md.EndCollapsable()

	// Google form
	md.AddLine("Have any feedback? Feel free to submit it")
	markup.Link("here", utils.GoogleFormURL)

	return md.Render()
}
