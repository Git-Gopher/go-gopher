package commands

import (
	"github.com/Git-Gopher/go-gopher/assess/options"
	log "github.com/sirupsen/logrus"
)

func LoadOptions(log *log.Logger) *options.Options {
	o := options.Options{}
	r := options.NewFileReader(log, &o)

	if err := r.Read("options.yml"); err != nil {
		log.Fatalf("failed to read options: %v", err)
	}

	return &o
}
