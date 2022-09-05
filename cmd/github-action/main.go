package main

import (
	"os"

	"github.com/Git-Gopher/go-gopher/commands"
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
	app.Action = commands.ActionCommand.Action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
