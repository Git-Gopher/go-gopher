package cli

import (
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "go-gopher-cli"
	app.HelpName = "go-gopher-cli"
	app.Usage = "A tool for analyzing GitHub repositories"
	app.Version = version.BuildVersion()
	app.Commands = []*cli.Command{
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
			Action: DownloadCommand,
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
			Action: AnalyzeCommand,
			Subcommands: []*cli.Command{
				{
					Name:   "url",
					Usage:  "analyze github project url",
					Action: AnalyzeUrlCommand,
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
					Action: AnalyzeBatchCommand,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
