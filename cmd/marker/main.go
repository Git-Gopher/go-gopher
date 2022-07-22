package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "go-gopher-marker"
	app.HelpName = "go-gopher-marker"
	app.Usage = "A tool for mark GitHub projects"
	app.Commands = []*cli.Command{
		{
			Name:  "single",
			Usage: "grade a single repository",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Usage:    "the url of the repository to grade",
					Aliases:  []string{"u"},
					Required: true,
				},
			},
			Action: singleCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}