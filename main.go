package main

import (
	"errors"
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/urfave/cli/v2"
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
			Action: func(c *cli.Context) error {
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
				current := cache.NewCache(gitModel)
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

				ghwf := workflow.GithubFlowWorkflow()
				violated, count, total, violations, err := ghwf.Analyze(enrichedModel, current, caches)
				if err != nil {
					log.Fatalf("Failed to analyze: %v\n", err)
				}

				render(violated, count, total, violations)

				return nil
			},
		},
		{
			Name:    "memory",
			Aliases: []string{"a"},
			Usage:   "detect a workflow for a given git project url",
			Action: func(c *cli.Context) error {
				url := c.Args().Get(0)

				repo, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					URL: url,
				})

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
				current := cache.NewCache(gitModel)
				caches, err := cache.ReadCaches()
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Cache file does not exist: %v", err)
				} else {
					log.Fatalf("Failed to load caches: %v", err)
				}

				ghwf := workflow.GithubFlowWorkflow()
				violated, count, total, violations, err := ghwf.Analyze(enrichedModel, current, caches)
				if err != nil {
					log.Printf("err: %v\n", err)
				}
				render(violated, count, total, violations)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func render(v, c, t int, vs []violation.Violation) {
	log.Printf("violated: %d\n", v)
	log.Printf("count: %d\n", c)
	log.Printf("total: %d\n", t)
	log.Printf("\n###### Violations ######\n")
	for _, violation := range vs {
		log.Println(violation.Display())
	}
}
