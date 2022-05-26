package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:  "analyze",
			Usage: "detect a workflow for a given git project url",
			Action: func(c *cli.Context) error {
				url := c.Args().Get(0)

				// Default to project repository
				if url == "" {
					url = "https://github.com/Git-Gopher/go-gopher"
				}

				r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					URL: url,
				})

				fmt.Println(err)
				// ... retrieves the branch pointed by HEAD
				ref, err := r.Head()
				fmt.Println(ref)

				// ... retrieves the commit history
				cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
				fmt.Println(cIter)

				// ... just iterates over the commits, printing it
				err = cIter.ForEach(func(c *object.Commit) error {
					fmt.Println(c)
					return nil
				})
				return nil

			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
