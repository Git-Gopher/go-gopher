package main

import (
	"context"
	"fmt"

	"github.com/Git-Gopher/go-gopher/config"
	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type githubConfig struct {
	GithubActions         bool   `env:"GITHUB_ACTIONS,required"`          // always set to true in GitHub Actions.
	GithubAction          string `env:"GITHUB_ACTION,required"`           // name of the action.
	GithubToken           string `env:"GITHUB_TOKEN,required"`            // token scoped to the repo.
	GithubWorkspace       string `env:"GITHUB_WORKSPACE,required"`        // path to the workspace.
	GithubRepository      string `env:"GITHUB_REPOSITORY,required"`       // example: octocat/Hello-World.
	GithubRepositoryOwner string `env:"GITHUB_REPOSITORY_OWNER,required"` // example: octocat.
	GithubActor           string `env:"GITHUB_ACTOR,required"`            // example: octocat.
	// Comma separated list of disallowed (blacklisted) users. Primarily used in the case of filtering bot accounts.
	LoginDisallowList string `env:"LOGIN_DISALLOW_LIST,default=github-classroom[bot]"`
}

// Load the environment variables from GitHub Actions.
func loadEnv(ctx context.Context) (*githubConfig, error) {
	var c githubConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	return &c, nil
}

// Fetch custom or default config. Fatal on bad custom config.
func readConfig(ctx *cli.Context) *config.Config {
	var cfg *config.Config
	var err error

	// Custom config
	if ctx.String("config") != "" {
		cfg, err = config.Read(ctx.String("config"))
		if err != nil {
			log.Fatalf("Failed to read custom config: %v", err)
		}
	} else {
		// Use default config
		cfg, err = config.Default()
		if err != nil {
			log.Fatalf("Failed to read default config: %v", err)
		}
	}

	return cfg
}
