package main

import (
	"os"

	"github.com/Git-Gopher/go-gopher/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})

	var cmd Commands = &Cmds{}

	app := cli.NewApp()

	app.Name = "go-gopher-workflow"
	app.HelpName = "go-gopher-workflow"
	app.Usage = "A tool for detecting workflows in GitHub repositories"
	app.Version = version.BuildVersion()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:      "url",
			Aliases:   []string{"u"},
			Category:  "Repository",
			Usage:     "grade a single repository with GitHub URL",
			UsageText: "go-gopher-workflow url <url> - evaluate repository with GitHub remote URL",
			ArgsUsage: "<url>",
			Action:    LoadFlags(cmd.SingleUrlCommand),
		},
		{
			Name:     "batch",
			Aliases:  []string{"b"},
			Category: "Repository",
			Usage:    "grade a set of repositories from a json file",
			UsageText: `go-gopher-workflow batch <repos.json> - 
					parse file with GitHub remote URL that has been generated from go-gopher-cli query`,
			ArgsUsage: "<repos.json>",
			Action:    LoadFlags(cmd.BatchUrlCommand),
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t", "gh"},
			Usage:   "GitHub token to access private repositories",
		},
		&cli.StringFlag{
			Name:        "env",
			Aliases:     []string{"e"},
			DefaultText: ".env",
			Usage:       "Environment file location. Default: .env",
		},
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "timeout in seconds before the repository is skipped",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
