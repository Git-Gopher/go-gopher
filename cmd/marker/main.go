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

	app.Name = "go-gopher-marker"
	app.HelpName = "go-gopher-marker"
	app.Usage = "A tool for mark GitHub projects"
	app.Version = version.BuildVersion()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:      "url",
			Aliases:   []string{"u"},
			Category:  "Marker",
			Usage:     "grade a single repository with GitHub URL",
			UsageText: "go-gopher-marker url <url> - grade repository with GitHub URL",
			ArgsUsage: "<url>",
			Action:    LoadFlags(cmd.SingleUrlCommand),
		},
		{
			Name:      "local",
			Aliases:   []string{"l"},
			Category:  "Marker",
			Usage:     "grade a single local repository",
			UsageText: "go-gopher-marker local <path> - grade local repository",
			ArgsUsage: "<folder>",
			Action:    LoadFlags(cmd.SingleLocalCommand),
		},
		{
			Name:      "folder",
			Aliases:   []string{"f"},
			Category:  "Marker",
			Usage:     "grade a folder of repositories",
			UsageText: "go-gopher-marker folder <path> - grade folder of repositories",
			ArgsUsage: "<folder>",
			Action:    LoadFlags(cmd.FolderLocalCommand),
		},
		{
			Name:     "generate",
			Category: "Utils",
			Usage:    "generate and reset configuration files. options.yml and .env files",
			Action:   SkipFlags(cmd.GenerateConfigCommand),
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t", "gh"},
			Usage:   "GitHub token to access private repositories",
		},
		&cli.StringFlag{
			Name:        "options",
			Aliases:     []string{"o", "opt"},
			DefaultText: "options.yml",
			Usage:       "Options file location. Default: options.yml",
		},
		&cli.StringFlag{
			Name:        "env",
			Aliases:     []string{"e"},
			DefaultText: ".env",
			Usage:       "Environment file location. Default: .env",
		},
		&cli.StringFlag{
			Name:        "lookup-path",
			DefaultText: "./data/se206-2022-beta-students.csv",
			Usage:       "student lookup csv file location. Default: ./data/se206-2022-beta-students.csv",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
