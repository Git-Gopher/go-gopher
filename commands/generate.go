package commands

import (
	"errors"
	"fmt"

	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/urfave/cli/v2"
)

func (c *Cmds) GenerateConfigCommand(cCtx *cli.Context, flags *Flags) error {
	r := options.NewFileReader(log.StandardLogger(), nil)
	if err := r.GenerateDefault(flags.OptionsDir); err != nil {
		if !errors.Is(err, utils.ErrSkipped) {
			return fmt.Errorf("failed to generate default options: %w", err)
		}
	}

	if err := utils.GenerateEnv(flags.EnvDir); err != nil {
		if !errors.Is(err, utils.ErrSkipped) {
			return fmt.Errorf("failed to generate env: %w", err)
		}
	}

	return nil
}
