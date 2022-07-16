package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "go-gopher-action"
	app.HelpName = "go-gopher-action"
	app.Usage = "A github action for analyzing GitHub repositories"
	app.Action = actionCommand

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
