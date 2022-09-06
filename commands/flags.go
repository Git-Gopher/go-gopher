package commands

import (
	"fmt"
	"os"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var errGitHubToken = fmt.Errorf("missing GitHub token, use --token or GITHUB_TOKEN")

type Flags struct {
	GithubToken string
	OptionsDir  string
	EnvDir      string
}

func NewFlags() *Flags {
	return &Flags{
		GithubToken: "",
		OptionsDir:  "options.yml",
		EnvDir:      ".env",
	}
}

type ActionWithFlagFunc func(cCtx *cli.Context, flags *Flags) error

func LoadFlags(command ActionWithFlagFunc) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		log.Infof("BuildVersion: %v\n", version.BuildVersion())

		flags := NewFlags()

		if cCtx.String("env") != "" {
			flags.EnvDir = cCtx.String("env")
		}

		_ = utils.Env(flags.EnvDir)

		if cCtx.String("options") != "" {
			flags.OptionsDir = cCtx.String("options")
		}

		if cCtx.String("token") != "" {
			flags.GithubToken = cCtx.String("token")
		} else {
			flags.GithubToken = os.Getenv("GITHUB_TOKEN")
		}

		if flags.GithubToken == "" {
			return errGitHubToken
		}

		_ = os.Setenv("GITHUB_TOKEN", flags.GithubToken)

		return command(cCtx, flags)
	}
}

func SkipFlags(command ActionWithFlagFunc) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		flags := NewFlags()

		if cCtx.String("env") != "" {
			flags.EnvDir = cCtx.String("env")
		}

		return command(cCtx, flags)
	}
}
