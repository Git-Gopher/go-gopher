package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/config"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

//nolint
func main() {
	utils.Environment(".env")
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "action",
			Aliases: []string{"a"},
			Usage:   "detect a workflow for current root",
			Action: func(ctx *cli.Context) error {
				// repository := os.Getenv("GITHUB_REPOSITORY")
				workspace := os.Getenv("GITHUB_WORKSPACE")
				if workspace == "" {
					log.Fatalf("GITHUB_WORKSPACE is not set")
				}
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
					var remotes []*git.Remote
					remotes, err = repo.Remotes()
					if err != nil {
						log.Fatalf("Could not get git repository remotes: %v\n", err)
					}

					if len(remotes) == 0 {
						log.Fatalf("No remotes present: %v\n", err)
					}

					// XXX: Use the first remote, assuming origin.
					urls := remotes[0].Config().URLs
					if len(urls) == 0 {
						log.Fatalf("No URLs present: %v\n", err)
					}

					url = urls[0]
				}
				owner, name, err := utils.OwnerNameFromUrl(url)
				if err != nil {
					log.Fatalf("Could not get owner and name from URL: %v\n", err)
				}

				githubModel, err := github.ScrapeGithubModel(owner, name)
				if err != nil {
					log.Fatalf("Could not create GithubModel: %v\n", err)
				}

				enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

				// Cache
				current := cache.NewCache(enrichedModel)
				caches, err := cache.ReadCaches()

				//nolint
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Cache file does not exist: %v", err)
					// Write a cache for current so that next run can use it
					if err = cache.WriteCaches([]*cache.Cache{current}); err != nil {
						log.Fatalf("Could not write cache: %v\n", err)
					}
				} else if err != nil {
					log.Fatalf("Failed to load caches: %v", err)
				} else {
				}

				cfg := readConfig(ctx)
				ghwf := workflow.GithubFlowWorkflow(cfg)
				violated, count, total, violations, err := ghwf.Analyze(enrichedModel, current, caches)
				if err != nil {
					log.Fatalf("Failed to analyze: %v\n", err)
				}

				workflowLog(violated, count, total, violations)

				// Set action outputs to a markdown summary.
				md := markup.NewMarkdown()
				md.
					Title("Workflow Summary").
					Collapsible("Violations", markup.NewMarkdown().Text("Stub!")).
					Collapsible("Suggestions", markup.NewMarkdown().Text("Stub!")).
					Collapsible("Authors", markup.NewMarkdown().Text("Stub!"))

				markup.Outputs("pr_summary", md.String())

				ghwf.Csv(workflow.DefaultCsvPath)

				return nil
			},
		},
		{
			Name:    "memory",
			Aliases: []string{"m"},
			Usage:   "detect a workflow for a given git project url",
			Action: func(ctx *cli.Context) error {
				url := ctx.Args().Get(0)

				repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					URL: url,
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

				githubModel, err := github.ScrapeGithubModel(owner, name)
				if err != nil {
					log.Fatalf("Could not scrape GithubModel: %v\n", err)
				}
				enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

				// Cache
				current := cache.NewCache(enrichedModel)
				caches, err := cache.ReadCaches()
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Cache file does not exist: %v", err)
				} else {
					log.Fatalf("Failed to load caches: %v", err)
				}

				cfg := readConfig(ctx)
				ghwf := workflow.GithubFlowWorkflow(cfg)
				violated, count, total, violations, err := ghwf.Analyze(enrichedModel, current, caches)
				if err != nil {
					log.Printf("err: %v\n", err)
				}
				workflowLog(violated, count, total, violations)

				return nil
			},
		},
		{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "debug/development subcommands",
			Subcommands: []*cli.Command{
				{
					Name:    "feature",
					Aliases: []string{"feat", "features"},
					Usage:   "detect a feature branching",
					// Example: `go-gopher feature https://github.com/Git-Gopher/tests test/two-parents-merged/0`
					Action: func(ctx *cli.Context) error {
						url := ctx.Args().Get(0)
						branch := ctx.Args().Get(1)

						if err := godotenv.Load(".env"); err != nil {
							log.Println("Error loading .env file")
						}

						token := os.Getenv("GITHUB_TOKEN")

						var branchRef plumbing.ReferenceName
						if branch != "" {
							branchRef = plumbing.NewBranchReferenceName(branch)
						}

						repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
							Auth: &http.BasicAuth{
								Username: "non-empty",
								Password: token,
							},
							URL:           url,
							ReferenceName: branchRef,
						})
						if err != nil {
							log.Fatalf("Failed to clone repository: %v", err)
						}

						gitModel, err := local.NewGitModel(repo)
						if err != nil {
							log.Fatalf("Could not create GitModel: %v\n", err)
						}

						enrichedModel := enriched.NewEnrichedModel(*gitModel, github.GithubModel{})

						d := detector.NewFeatureBranchDetector()
						if err := d.Run(enrichedModel); err != nil {
							log.Fatalf("Failed to run weighted detectors: %v", err)
						}
						v, co, t, vs := d.Result()

						fmt.Printf("\n## Detector Type: %T ##\n", d)
						workflowLog(v, co, t, vs)

						return nil
					},
				},
				{
					Name:    "diff-distance",
					Aliases: []string{"dd"},
					Usage:   "detect diff distance",
					// Example: `go-gopher dd https://github.com/Git-Gopher/tests test/two-parents-merged/0`
					Action: func(ctx *cli.Context) error {
						url := ctx.Args().Get(0)
						branch := ctx.Args().Get(1)

						if err := godotenv.Load(".env"); err != nil {
							log.Println("Error loading .env file")
						}

						token := os.Getenv("GITHUB_TOKEN")

						var branchRef plumbing.ReferenceName
						if branch != "" {
							branchRef = plumbing.NewBranchReferenceName(branch)
						}

						repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
							Auth: &http.BasicAuth{
								Username: "non-empty",
								Password: token,
							},
							URL:           url,
							ReferenceName: branchRef,
						})
						if err != nil {
							log.Fatalf("Failed to clone repository: %v", err)
						}

						gitModel, err := local.NewGitModel(repo)
						if err != nil {
							log.Fatalf("Could not create GitModel: %v\n", err)
						}

						enrichedModel := enriched.NewEnrichedModel(*gitModel, github.GithubModel{})

						d := detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation())
						if err := d.Run(enrichedModel); err != nil {
							log.Fatalf("Failed to run weighted detectors: %v", err)
						}
						v, co, t, vs := d.Result()

						fmt.Printf("\n## Detector Type: %T ##\n", d)
						workflowLog(v, co, t, vs)

						return nil
					},
				},
			},
		},
		{
			Name:    "local",
			Aliases: []string{"m"},
			Usage:   "detect a workflow for a given git project url",
			// Example: `go-gopher local https://github.com/Git-Gopher/tests test/two-parents-merged/0`
			Action: func(ctx *cli.Context) error {
				url := ctx.Args().Get(0)
				branch := ctx.Args().Get(1)

				if err := godotenv.Load(".env"); err != nil {
					log.Println("Error loading .env file")
				}

				token := os.Getenv("GITHUB_TOKEN")

				var branchRef plumbing.ReferenceName
				if branch != "" {
					branchRef = plumbing.NewBranchReferenceName(branch)
				}

				repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					Auth: &http.BasicAuth{
						Username: "non-empty",
						Password: token,
					},
					URL:           url,
					ReferenceName: branchRef,
				})
				if err != nil {
					log.Fatalf("Failed to clone repository: %v", err)
				}

				gitModel, err := local.NewGitModel(repo)
				if err != nil {
					log.Fatalf("Could not create GitModel: %v\n", err)
				}

				enrichedModel := enriched.NewEnrichedModel(*gitModel, github.GithubModel{})

				for _, detector := range workflow.LocalDetectors() {
					if err := detector.Run(enrichedModel); err != nil {
						log.Fatalf("Failed to run weighted detectors: %v", err)
					}
					v, c, t, vs := detector.Result()

					fmt.Printf("\n## Detector Type: %T ##\n", detector)
					workflowLog(v, c, t, vs)
				}

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
				&cli.StringFlag{
					Name:     "logging",
					Usage:    "enable logging, output to the set file name",
					Aliases:  []string{"l"},
					Required: false,
				},
			},

			Subcommands: []*cli.Command{
				{
					Name:  "url",
					Usage: "remove an existing template",
					Action: func(ctx *cli.Context) error {
						url := ctx.Args().Get(0)

						repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
							URL: url,
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

						githubModel, err := github.ScrapeGithubModel(owner, name)
						if err != nil {
							log.Fatalf("Could not scrape GithubModel: %v\n", err)
						}
						enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)
						fmt.Printf("enrichedModel: %v\n", enrichedModel)

						return nil
					},
				},
				{
					Name:  "local",
					Usage: "add a new template",
					Action: func(ctx *cli.Context) error {
						fmt.Println("new task template: ", ctx.Args().First())
						return nil
					},
				},
				{
					Name:  "local-batch",
					Usage: "batch",
					Action: func(ctx *cli.Context) error {
						fmt.Println("new task template: ", ctx.Args().First())
						return nil
					},
				},
			},
			Action: func(ctx *cli.Context) error {
				fmt.Println("detect path or url and apply batch, local or url to it")

				logger, _ := zap.NewProduction()
				zap.ReplaceGlobals(logger)

				defer logger.Sync() // flushes buffer, if any
				sugar := logger.Sugar()
				url := "asdf"
				sugar.Infof("Failed to fetch URL: %s", url)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Analyze a model
func analyze() {
}

// Print violation summary to IO, Split by severity with author association.
func workflowLog(v, c, t int, vs []violation.Violation) {
	var violations, suggestions []violation.Violation
	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			violations = append(violations, v)
		case violation.Suggestion:
			suggestions = append(suggestions, v)
		}
	}

	var vsd string
	for _, v := range violations {
		vsd += v.Display()
	}
	markup.Group("Violations", vsd)

	var ssd string
	for _, v := range suggestions {
		ssd += v.Display()
	}
	markup.Group("Suggestions", ssd)

	var asd string
	authors := make(map[string]int)
	for _, v := range vs {
		a, err := v.Author()
		if err != nil {
			continue
		}
		authors[a.Login]++
	}

	for author, count := range authors {
		asd += fmt.Sprintf("%s: %d\n", author, count)
	}

	asd += fmt.Sprintf("violated: %d\n", v)
	asd += fmt.Sprintf("count: %d\n", c)
	asd += fmt.Sprintf("total: %d\n", t)
	markup.Group("Summary", asd)
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
	} else
	// Use default config
	{
		cfg, err = config.Default()
		if err != nil {
			log.Fatalf("Failed to read default config: %v", err)
		}
	}

	return cfg
}

func logging(ctx *cli.Context) {
}
