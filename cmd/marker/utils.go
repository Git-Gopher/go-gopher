package main

import (
	"github.com/Git-Gopher/go-gopher/assess/options"
	log "github.com/sirupsen/logrus"
)

func LoadOptions(log *log.Logger) *options.Options {
	o := options.Options{}
	r := options.NewFileReader(log, &o)
	r.Read("options.yml")

	return &o
}
