package main

import (
	"errors"
	"os"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	errGitHubToken = errors.New("missing GitHub token, use --token or GITHUB_TOKEN")
	errTimeout     = errors.New("invalid timeout value")
)

type Flags struct {
	GithubToken string
	EnvDir      string
	Timeout     int
}

func NewFlags() *Flags {
	return &Flags{
		GithubToken: "",
		EnvDir:      ".env",
		Timeout:     3 * 60, // 3 minutes
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

		if cCtx.String("token") != "" {
			flags.GithubToken = cCtx.String("token")
		} else {
			flags.GithubToken = os.Getenv("GITHUB_TOKEN")
		}

		if flags.GithubToken == "" {
			return errGitHubToken
		}

		_ = os.Setenv("GITHUB_TOKEN", flags.GithubToken)

		if cCtx.Int("timeout") == 0 {
			flags.Timeout = 3 * 60 // 3 minutes
		} else if cCtx.Int("timeout") < 0 {
			return errTimeout
		} else {
			flags.Timeout = cCtx.Int("timeout")
		}

		return command(cCtx, flags)
	}
}
