package commands

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/google/go-github/v45/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var DownloadCommand = &cli.Command{
	Name:  "download",
	Usage: "Download artifact logs from workflow runs for a batch of repositories",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "output directory",
			DefaultText: "output",
			Value:       "output",
			Required:    false,
		},
	},
	Action: func(ctx *cli.Context) error {
		utils.Environment(".env")
		org := ctx.Args().Get(0)
		out := ctx.String("output")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx.Context, ts)
		client := github.NewClient(tc)

		var logs []interface{}
		err := os.MkdirAll(out, 0o750)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		log.Printf("Fetching repositories for organization %s...\n", org)
		repos, _, err := client.Repositories.ListByOrg(ctx.Context, org, nil)
		if err != nil {
			log.Fatalf("Could not fetch orginisation repositories: %v", err)
		}

		for _, r := range repos {
			var arts *github.ArtifactList
			arts, _, err = client.Actions.ListArtifacts(ctx.Context, org, *r.Name, nil)
			if err != nil {
				log.Fatalf("Could not fetch artifact list: %v", err)
			}

			for _, a := range arts.Artifacts {
				log.Printf("Downloading artifacts for %s/%s...\n", org, *r.Name)
				var url *url.URL
				url, _, err = client.Actions.DownloadArtifact(ctx.Context, org, *r.Name, *a.ID, true)
				if err != nil {
					log.Fatalf("could not fetch artifact url: %v", err)
				}

				pathZip := fmt.Sprintf("%s/log-%s-%s-%d.zip", out, org, *r.Name, *a.ID)
				pathJson := fmt.Sprintf("%s/log-%s-%s-%d", out, org, *r.Name, *a.ID)

				log.Printf("Downloading artifact %s...\n", pathZip)
				err = utils.DownloadFile(pathZip, url.String())
				if err != nil {
					log.Fatalf("could not download artifact: %v", err)
				}

				log.Printf("Unzipping %s...\n", pathZip)
				if err = utils.Unzip(pathZip, pathJson); err != nil {
					log.Fatalf("failed to unzip log contents: %v", err)
				}
			}
		}

		logCount := 0
		if err = filepath.Walk(out, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			if matched, err := filepath.Match("log-go-gopher*.json", filepath.Base(path)); err != nil {
				return fmt.Errorf("could not match log file: %w", err)
			} else if matched {
				logCount++

				log.Printf("Appending %s to merged log...", path)
				file, err := os.ReadFile(filepath.Clean(path))
				if err != nil {
					log.Fatalf("could not read log file: %v", err)
				}

				var data interface{}
				if err = json.Unmarshal(file, &data); err != nil {
					log.Fatalf("failed to unmarshal log: %v", err)
				}

				logs = append(logs, data)
			}

			return nil
		}); err != nil {
			log.Fatalf("failed to walk log directory: %v", err)
		}

		bytes, err := json.MarshalIndent(logs, "", " ")
		if err != nil {
			return fmt.Errorf("error marshaling merged log: %w", err)
		}

		logPath := fmt.Sprintf("%s/merged-log-%s.json", out, org)
		if err := os.WriteFile(logPath, bytes, 0o600); err != nil {
			return fmt.Errorf("error writing merged log: %w", err)
		}

		log.Printf("Downloaded %d logs from %s and merged to %s", logCount, "Git-Gopher", logPath)

		return nil
	},
}
