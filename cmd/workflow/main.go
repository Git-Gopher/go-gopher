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
			Name:      "batch",
			Aliases:   []string{"b"},
			Category:  "Repository",
			Usage:     "grade a set of repositories from a json file",
			UsageText: "go-gopher-workflow batch <file.json> - parse file with GitHub remote URL",
			ArgsUsage: "<url>",
			Action:    LoadFlags(cmd.SingleUrlCommand),
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
		&cli.StringFlag{
			Name:  "json",
			Usage: "Use json file to load workflow urls",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
