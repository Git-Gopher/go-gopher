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

	app := cli.NewApp()

	app.Name = "go-gopher-marker"
	app.HelpName = "go-gopher-marker"
	app.Usage = "A tool for mark GitHub projects"
	app.Version = version.BuildVersion()
	app.Commands = []*cli.Command{
		{
			Name:   "url",
			Usage:  "grade a single repository with GitHub URL",
			Action: singleUrlCommand,
		},
		{
			Name:   "local",
			Usage:  "grade a single local repository",
			Action: singleLocalCommand,
		},
		{
			Name:   "folder",
			Usage:  "grade a folder of repositories",
			Action: folderLocalCommand,
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "token",
			Usage: "GitHub token",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
