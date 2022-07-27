package main

import (
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "go-gopher"
	app.HelpName = "go-gopher"
	app.Usage = "A tool for analyzing GitHub repositories"
	app.Version = version.BuildVersion()
	// app.Action = actionCommand.

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
