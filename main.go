package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "analyze",
			Aliases: []string{"a"},
			Usage:   "detect a workflow for a given git project url",
			Action: func(c *cli.Context) error {
				url := c.Args().Get(0)

				// Default to project repository
				if url == "" {
					url = "https://github.com/jgm/pandoc"
				}

				r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					URL: url,
				})

				// ... retrieves the branch pointed by HEAD
				ref, _ := r.Head()

				log.Println(ref)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
