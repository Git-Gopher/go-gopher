package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/discord"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var errOwnerMismatch = errors.New("owner mismatch")

func actionCommand(cCtx *cli.Context) error {
	log.Printf("BuildVersion: %s", version.BuildVersion())
	// Load the environment variables from GitHub Actions.
	config, err := loadEnv(cCtx.Context)
	if err != nil {
		return fmt.Errorf("failed to load env: %w", err)
	}

	// Open the repository.
	repo, err := git.PlainOpen(config.GithubWorkspace)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}

	// GithubURL fallback.
	githubURL, err := utils.Url(repo)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}

	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return fmt.Errorf("failed to get owner and repo name: %w", err)
	}
	if config.GithubRepositoryOwner != repoOwner {
		return fmt.Errorf("%w: %s != %s", errOwnerMismatch, repoOwner, config.GithubRepositoryOwner)
	}

	// Create enrichedModel.
	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("failed to create enriched model: %w", err)
	}

	// Create cache.
	current := cache.NewCache(enrichedModel)

	// Populate authors from enrichedModel.
	authors := enriched.PopulateAuthors(enrichedModel)

	// Read cache.
	caches, err := cache.Read()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read caches: %w", err)
		}

		// Write a cache for current so that next run can use it.
		if err = cache.Write([]*cache.Cache{current}); err != nil {
			return fmt.Errorf("failed to write cache: %w", err)
		}
	}

	cfg := readConfig(cCtx)
	ghwf := workflow.GithubFlowWorkflow(cfg)
	violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
	if err != nil {
		log.Fatalf("Failed to analyze: %v\n", err)
	}

	if config.LoginWhiteList != "" {
		whitelist := strings.Split(config.LoginWhiteList, ",")
		violations = violation.FilterByLogin(violations, whitelist)
	}

	workflowSummary(authors, violated, count, total, violations)

	summary := markdownSummary(authors, violations)
	markup.Outputs("pr_summary", summary)

	fn, err := ghwf.WriteLog(*enrichedModel, cfg)
	if err != nil {
		log.Printf("Could not write json log: %v", err)

		return nil
	}

	if err = discord.SendLog(fn); err != nil {
		log.Printf("Could not write json log to discord: %v", err)

		return nil
	}

	return nil
}

// Print violation summary to IO, Split by severity with author association.
func workflowSummary(authors utils.Authors, v, c, t int, vs []violation.Violation) {
	var violations, suggestions []violation.Violation
	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			violations = append(violations, v)
		case violation.Suggestion:
			suggestions = append(suggestions, v)
		}
	}

	var vSd strings.Builder
	for _, v := range violations {
		vSd.WriteString(v.Display(authors))
	}
	markup.Group("Violations", vSd.String())

	var sSd strings.Builder
	for _, v := range suggestions {
		sSd.WriteString(v.Display(authors))
	}
	markup.Group("Suggestions", sSd.String())

	var aSd strings.Builder
	counts := make(map[string]int)
	for _, v := range vs {
		email := v.Email()
		login, err := authors.Find(email)
		if err != nil {
			continue
		}
		counts[*login]++
	}

	for login, count := range counts {
		aSd.WriteString(fmt.Sprintf("%s: %d\n", login, count))
	}

	aSd.WriteString(fmt.Sprintf("violated: %d\n", v))
	aSd.WriteString(fmt.Sprintf("count: %d\n", c))
	aSd.WriteString(fmt.Sprintf("total: %d\n", t))
	markup.Group("Summary", aSd.String())
}

// Helper function to create a markdown summary of the violations.
// nolint: gocognit
func markdownSummary(authors utils.Authors, vs []violation.Violation) string {
	md := markup.CreateMarkdown("Workflow Summary")
	md.AddLine(fmt.Sprintf("Created with git-gopher version `%s`", version.BuildVersion()))

	// Separate violation types.
	var violations []violation.Violation
	var suggestions []violation.Violation

	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			if v.Current() {
				violations = append(violations, v)
			}
		case violation.Suggestion:
			if v.Current() {
				suggestions = append(suggestions, v)
			}
		default:
			log.Printf("Unknown violation severity: %v", v.Severity())
		}
	}

	if len(violations) > 0 {
		headers := []string{"Violation", "Message", "Advice", "Author"}
		rows := make([][]string, len(violations))

		for i, v := range violations {
			row := make([]string, len(headers))
			name := v.Name()
			row[0] = name
			message := v.Message()
			row[1] = message

			suggestion, err := v.Suggestion()
			if err != nil {
				suggestion = ""
			}
			row[2] = suggestion

			usernamePtr, err := authors.Find(v.Email())
			if err != nil || usernamePtr == nil {
				row[3] = "unknown"
			} else {
				row[3] = markup.Author(*usernamePtr).Markdown()
			}

			rows[i] = row
		}

		md.BeginCollapsable("Violations")
		md.Table(headers, rows)
		md.EndCollapsable()
	}

	if len(suggestions) > 0 {
		headers := []string{"Suggestion", "Message", "Advice", "Author"}
		rows := make([][]string, len(suggestions))

		for i, v := range suggestions {
			row := make([]string, len(headers))
			name := v.Name()
			row[0] = name
			message := v.Message()
			row[1] = message

			suggestion, err := v.Suggestion()
			if err != nil {
				suggestion = ""
			}
			row[2] = suggestion

			usernamePtr, err := authors.Find(v.Email())
			if err != nil || usernamePtr == nil {
				row[3] = "unknown"
			} else {
				row[3] = markup.Author(*usernamePtr).Markdown()
			}

			rows[i] = row
		}

		md.BeginCollapsable("Suggestions")
		md.Table(headers, rows)
		md.EndCollapsable()
	}

	workflowUrl := os.Getenv("WORKFLOW_URL")
	if (len(violations)+len(suggestions)) < len(vs) && workflowUrl != "" {
		md.AddLine(fmt.Sprintf(`There still exist some violations beyond the scope of this pull request, 
			please view the full log [here](%s)`, workflowUrl))
	}

	// Google form
	md.AddLine(fmt.Sprintf("Have any feedback? Feel free to submit it [here](%s)", utils.GoogleFormURL))

	return md.Render()
}
