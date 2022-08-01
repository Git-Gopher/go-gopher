package options

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed options.reference.yml
var optionsByte []byte

type FileReader struct {
	log     *log.Logger
	options *Options
}

func NewFileReader(log *log.Logger, options *Options) *FileReader {
	return &FileReader{
		log:     log,
		options: options,
	}
}

func (r *FileReader) Read(file string) error {
	r.log.Info("Reading options from file: ", file)

	if _, err := os.Stat(file); err != nil {
		r.log.Info("Generating default options")
		if err := r.generateDefault(file); err != nil {
			return fmt.Errorf("can't generate default options: %w", err)
		}
	}

	viper.SetConfigFile(file)

	// Assume YAML if the file has no extension.
	if filepath.Ext(file) == "" {
		viper.SetConfigType("yaml")
	}

	if err := r.parseOption(); err != nil {
		return fmt.Errorf("can't parse default options: %w", err)
	}

	return nil
}

func (r *FileReader) generateDefault(file string) error {
	err := ioutil.WriteFile(file, optionsByte, 0o600)
	if err != nil {
		return fmt.Errorf("can't write default options: %w", err)
	}

	return nil
}

func (r *FileReader) parseOption() error {
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			r.log.Info("No options file found")

			return nil
		}

		return fmt.Errorf("can't read viper config: %w", err)
	}

	if err := viper.Unmarshal(r.options); err != nil {
		return fmt.Errorf("can't unmarshal config by viper: %w", err)
	}

	return nil
}
