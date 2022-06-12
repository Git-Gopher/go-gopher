package main

import (
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
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
				// sha := os.Getenv("GITHUB_SHA") // commit sha triggered
				// ref := os.Getenv("GITHUB_REF") // branch ref triggered

				r, err := git.PlainOpen(workspace)
				if err != nil {
					log.Printf("cannot read repo: %v\n", err)
					os.Exit(1)
				}

				repo, err := local.NewGitModel(r)
				if err != nil {
					log.Printf("Could not create GitModel: %v\n", err)
					os.Exit(1)
				}
				ghwf := workflow.GithubFlowWorkflow()
				violated, count, total, err := ghwf.Analyze(repo)
				if err != nil {
					log.Printf("Failed to analyze: %v\n", err)
					os.Exit(1)
				}
				log.Printf("violated: %d\n", violated)
				log.Printf("count: %d\n", count)
				log.Printf("total: %d\n", total)

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
				ghwf := workflow.GithubFlowWorkflow()
				violated, count, total, err := ghwf.Analyze(repo)
				if err != nil {
					log.Printf("err: %v\n", err)
				}
				log.Printf("violated: %d\n", violated)
				log.Printf("count: %d\n", count)
				log.Printf("total: %d\n", total)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
