package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v45/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func AnalyzeCommand(ctx *cli.Context) error {
	utils.Environment(".env")
	workspace := utils.EnvGithubWorkspace()

	// repository := os.Getenv("GITHUB_REPOSITORY")
	// sha := os.Getenv("GITHUB_SHA") // commit sha triggered
	// ref := os.Getenv("GITHUB_REF") // branch ref triggered

	// Repo
	repo, err := git.PlainOpen(workspace)
	if err != nil {
		log.Fatalf("cannot read repo: %v\n", err)
	}

	gitModel, err := local.NewGitModel(repo)
	if err != nil {
		log.Fatalf("Could not create GitModel: %v\n", err)
	}

	url := os.Getenv("GITHUB_URL")
	if url == "" {
		url, err = utils.Url(repo)
		if err != nil {
			log.Fatalf("Could get url from repository: %v\n", err)
		}
		log.Printf("GITHUB_URL is not set, falling back to \"%s\"...\n", url)
	}
	owner, name, err := utils.OwnerNameFromUrl(url)
	if err != nil {
		log.Fatalf("Could not get owner and name from URL: %v\n", err)
	}

	remoteModel, err := remote.ScrapeRemoteModel(owner, name)
	if err != nil {
		log.Fatalf("Could not create RemoteModel: %v\n", err)
	}

	enrichedModel := enriched.NewEnrichedModel(*gitModel, *remoteModel)

	// Authors
	authors := enriched.PopulateAuthors(enrichedModel)

	// Cache
	current := cache.NewCache(enrichedModel)
	caches, err := cache.Read()

	//nolint
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Cache file does not exist: %v", err)
		// Write a cache for current so that next run can use it
		if err = cache.Write([]*cache.Cache{current}); err != nil {
			log.Fatalf("Could not write cache: %v\n", err)
		}
	} else if err != nil {
		log.Fatalf("Failed to load caches: %v", err)
	}

	cfg := utils.ReadConfig(ctx)
	ghwf := workflow.GithubFlowWorkflow(cfg)
	violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
	if err != nil {
		log.Fatalf("Failed to analyze: %v\n", err)
	}

	workflow.TextSummary(authors, violated, count, total, violations)

	// Set action outputs to a markdown summary.
	summary := workflow.MarkdownSummary(authors, violations)
	markup.Outputs("pr_summary", summary)

	err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
	if err != nil {
		log.Fatalf("Could not create csv summary: %v", err)
	}

	err = ghwf.WriteLog(*enrichedModel, cfg)
	if err != nil {
		log.Fatalf("Could not write json log: %v", err)
	}

	return nil
}

func DownloadCommand(ctx *cli.Context) error {
	utils.Environment(".env")
	org := ctx.Args().Get(0)
	out := ctx.String("output")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx.Context, ts)
	client := github.NewClient(tc)

	var logs []interface{}
	err := os.MkdirAll(out, 0o755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	log.Printf("Fetching repositories for organisation %s...\n", org)
	repos, _, err := client.Repositories.ListByOrg(ctx.Context, org, nil)
	if err != nil {
		log.Fatalf("Could not fetch orginisation repositories: %v", err)
	}

	for _, r := range repos {
		arts, _, err := client.Actions.ListArtifacts(ctx.Context, org, *r.Name, nil)
		if err != nil {
			log.Fatalf("Could not fetch artifact list: %v", err)
		}

		for _, a := range arts.Artifacts {
			log.Printf("Downloading artifacts for %s/%s...\n", org, *r.Name)
			url, _, err := client.Actions.DownloadArtifact(ctx.Context, org, *r.Name, *a.ID, true)
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
			utils.Unzip(pathZip, pathJson)
		}
	}

	logCount := 0
	filepath.Walk(out, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("log-go-gopher*.json", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			logCount++

			log.Printf("Appending %s to merged log...", path)
			file, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Could not read log file: %v", err)
			}

			var data interface{}
			json.Unmarshal(file, &data)
			logs = append(logs, data)
		}
		return nil
	})

	bytes, err := json.MarshalIndent(logs, "", " ")
	if err != nil {
		return fmt.Errorf("error marshaling merged log: %w", err)
	}

	logPath := fmt.Sprintf("%s/merged-log-%s.json", out, org)
	if err := ioutil.WriteFile(logPath, bytes, 0o600); err != nil {
		return fmt.Errorf("error writing merged log: %w", err)
	}

	log.Printf("Downloaded %d logs from %s and merged to %s", logCount, "Git-Gopher", logPath)

	return nil
}

func AnalyzeUrlCommand(ctx *cli.Context) error {
	utils.Environment(".env")
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	url := ctx.Args().Get(0)

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
		Auth: &githttp.BasicAuth{
			Username: "non-empty",
			Password: token,
		},
	})
	if err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	gitModel, err := local.NewGitModel(repo)
	if err != nil {
		log.Fatalf("Could not create GitModel: %v\n", err)
	}

	owner, name, err := utils.OwnerNameFromUrl(url)
	if err != nil {
		log.Fatalf("Could not get owner and name from URL: %v\n", err)
	}

	githubModel, err := remote.ScrapeRemoteModel(owner, name)
	if err != nil {
		log.Fatalf("Could not scrape GithubModel: %v\n", err)
	}
	enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

	// Authors
	authors := enriched.PopulateAuthors(enrichedModel)

	// Cache
	current := cache.NewCache(enrichedModel)
	caches, err := cache.Read()

	//nolint
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Cache file does not exist: %v", err)
		// Write a cache for current so that next run can use it
		if err = cache.Write([]*cache.Cache{current}); err != nil {
			log.Fatalf("Could not write cache: %v\n", err)
		}
	} else if err != nil {
		log.Fatalf("Failed to load caches: %v", err)
	}

	cfg := utils.ReadConfig(ctx)
	ghwf := workflow.GithubFlowWorkflow(cfg)
	violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
	if err != nil {
		log.Fatalf("Failed to analyze: %v\n", err)
	}

	workflow.TextSummary(authors, violated, count, total, violations)

	if ctx.Bool("csv") {
		err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
		if err != nil {
			log.Fatalf("Could not create csv summary: %v", err)
		}
	}

	if ctx.Bool("logging") {
		err = ghwf.WriteLog(*enrichedModel, cfg)
		if err != nil {
			log.Fatalf("Could not write json log: %v", err)
		}
	}

	return nil
}

func AnalyzeBatchCommand(ctx *cli.Context) error {
	utils.Environment(".env")
	path := ctx.Args().Get(0)
	if path == "" {
		path = "./"
		log.Printf("No path provided, using current directory (\"%s\") as target path", path)

	}

	var ps []string
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			ps = append(ps, filepath.Dir(path))
			return filepath.SkipDir
		}

		return nil
	})
	if len(ps) == 0 {
		log.Fatalf("Could not detect any git repositories within the directiory: \"%s\"", path)
	}

	cfg := utils.ReadConfig(ctx)
	ghwf := workflow.GithubFlowWorkflow(cfg)

	for _, p := range ps {
		repo, err := git.PlainOpen(p)
		if err != nil {
			log.Fatalf("Failed to clone repository: %v", err)
		}

		gitModel, err := local.NewGitModel(repo)
		if err != nil {
			log.Fatalf("Could not create GitModel: %v\n", err)
		}

		url, err := utils.Url(repo)
		if err != nil {
			log.Fatalf("Could get url from repository: \"%v\", does it have any remotes?", err)
		}

		owner, name, err := utils.OwnerNameFromUrl(url)
		if err != nil {
			log.Fatalf("Could get the owner and name from URL: %v", err)
		}

		githubModel, err := remote.ScrapeRemoteModel(owner, name)
		if err != nil {
			log.Fatalf("Could not create GithubModel: %v\n", err)
		}

		enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

		// Authors
		authors := enriched.PopulateAuthors(enrichedModel)

		v, c, t, vs, err := ghwf.Analyze(enrichedModel, authors, nil, nil)
		if err != nil {
			log.Fatalf("Failed to analyze: %v\n", err)
		}

		workflow.TextSummary(authors, v, c, t, vs)
		nameCsv := fmt.Sprintf("batch-%s.csv", filepath.Base(path))
		if ctx.Bool("csv") {
			err = ghwf.Csv(nameCsv, enrichedModel.Name, enrichedModel.URL)
			if err != nil {
				log.Fatalf("Could not create csv summary: %v", err)
			}
		}

		if ctx.Bool("logging") {
			err = ghwf.WriteLog(*enrichedModel, cfg)
			if err != nil {
				log.Fatalf("Could not write json log: %v", err)
			}
		}
	}
	return nil
}
