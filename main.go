package main

import (
	"errors"
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
	workflow "github.com/Git-Gopher/go-gopher/worflow"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/urfave/cli/v2"
)

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
				r, err := git.PlainOpen(workspace)
				if err != nil {
					log.Fatalf("cannot read repo: %v\n", err)
				}

				repo, err := local.NewGitModel(r)
				if err != nil {
					log.Fatalf("Could not create GitModel: %v\n", err)
				}

				// Cache
				current := cache.NewCache(repo)
				caches, err := cache.ReadCaches()
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Cache file does not exist: %v", err)
					// Write a cache for current so that next run can use it
					cache.WriteCaches([]*cache.Cache{current})
				} else if err != nil {
					log.Fatalf("Failed to load caches: %v", err)
				} else {
					ghwf := workflow.GithubFlowWorkflow()
					violated, count, total, violations, err := ghwf.Analyze(repo, current, caches)
					if err != nil {
						log.Fatalf("Failed to analyze: %v\n", err)
					}

					render(violated, count, total, violations)
					for _, violation := range violations {
						log.Println(violation.Message())
					}
				}

				return nil
			},
		},
		{
			Name:    "memory",
			Aliases: []string{"a"},
			Usage:   "detect a workflow for a given git project url",
			Action: func(c *cli.Context) error {
				url := c.Args().Get(0)

				r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					URL: url,
				})

				// ... retrieves the branch pointed by HEAD
				ref, _ := r.Head()

				log.Println(ref)

				// TODO:
				// From the url create a enriched model...
				repo, err := local.NewGitModel(r)
				if err != nil {
					log.Printf("err: %v\n", err)
				}

				// Cache
				current := cache.NewCache(repo)
				caches, err := cache.ReadCaches()
				if errors.Is(err, os.ErrNotExist) {
					caches = nil
					log.Printf("Cache file does not exist: %v", err)

				} else {
					log.Fatalf("Failed to load caches: %v", err)
				}

				ghwf := workflow.GithubFlowWorkflow()
				violated, count, total, violations, err := ghwf.Analyze(repo, current, caches)
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
		log.Println(violation.Message())
	}
}
