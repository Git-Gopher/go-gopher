package main

import (
	"os"

	"github.com/Git-Gopher/go-gopher/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "go-gopher-action"
	app.HelpName = "go-gopher-action"
	app.Usage = "A github action for analyzing GitHub repositories"
	app.Version = version.BuildVersion()
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "disable-pr-comment",
			Required: false,
			Usage:    "disable pr comment output which cleans up stdout log",
		},
	}
	app.Action = actionCommand

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
