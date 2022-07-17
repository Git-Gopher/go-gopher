package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "go-gopher"
	app.HelpName = "go-gopher"
	app.Usage = "A tool for analyzing GitHub repositories"
	// app.Action = actionCommand.

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
