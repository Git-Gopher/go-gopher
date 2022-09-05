package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/config"
	"github.com/Git-Gopher/go-gopher/discord"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v45/github"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var Developers = []string{"wqsz7xn", "scorpionknifes"}

//nolint:all
func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = version.BuildVersion()
	app.Commands = []*cli.Command{
		{
			Name:    "action",
			Aliases: []string{"a"},
			Usage:   "detect a workflow for current root",
			Action: func(ctx *cli.Context) error {
				utils.Environment(".env")
				workspace := utils.EnvGithubWorkspace()

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
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("Cache file does not exist: %v", err)
					// Write a cache for current so that next run can use it
					if err = cache.Write([]*cache.Cache{current}); err != nil {
						log.Fatalf("Could not write cache: %v\n", err)
					}
				} else if err != nil {
					log.Fatalf("Failed to load caches: %v", err)
				}

				cfg := readConfig(ctx)
				ghwf := workflow.GithubFlowWorkflow(cfg)
				violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
				if err != nil {
					log.Fatalf("Failed to analyze: %v\n", err)
				}

				workflowSummary(authors, violated, count, total, violations)

				// Set action outputs to a markdown summary.
				summary := markdownSummary(authors, violations)
				markup.Outputs("pr_summary", summary)

				err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
				if err != nil {
					log.Fatalf("Could not create csv summary: %v", err)
				}

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
			},
		},
		{
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
		},
		{
			Name:  "team",
			Usage: "Add wqsz7xn and scorpionknifes as team members to each repository within the organization",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "prefix",
					Aliases:  []string{"p"},
					Usage:    "repository prefix",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "token",
					Aliases:  []string{"t"},
					Usage:    "github token",
					Required: false,
				},
			},
			Action: func(ctx *cli.Context) error {
				utils.Environment(".env")
				organizationName := ctx.Args().Get(0)
				prefix := ctx.String("prefix")

				if prefix == "" {
					log.Printf("No repository prefix set via flat, all repositories within organization will have team added...")
				} else {
					log.Printf("Using repository prefix of %s...", prefix)
				}

				token := ctx.String("token")
				if token == "" {
					log.Printf("No github token passed in via flag, using environment file instead...")
					token = os.Getenv("GITHUB_TOKEN")
					if token == "" {
						log.Fatalf("Unable to find github token from flag or from environment file, exiting...")
					}
				}

				ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
				tc := oauth2.NewClient(ctx.Context, ts)
				client := github.NewClient(tc)

				log.Printf("Fetching organization %s...\n", organizationName)
				organization, _, err := client.Organizations.Get(ctx.Context, organizationName)
				if err != nil {
					log.Fatalf("Could not fetch orginization: %v", err)
				}

				teamSlug := "git-gopher"
				teamPermission := "pull"

				team, res, err := client.Teams.GetTeamBySlug(ctx.Context, organizationName, teamSlug)
				if err != nil && res.StatusCode != 404 {
					fmt.Printf("res.StatusCode: %v\n", res.StatusCode)
					log.Fatalf("Failed to fetch team by slug: %s, %v", res.Status, err)
				}

				// Team exists, add team to the rest of the organizations that don't have the team added.
				if team != nil {
					log.Printf(`Team %s for organization %s already exists.
						Adding team to all new repositories (duplicates don't matter)...`, teamSlug, organizationName)

					filteredRepos, err := fetchAllRepositoriesByPrefix(ctx, client, organizationName, prefix)
					if err != nil {
						log.Fatalf("Failed to fetch all repositories by prefix: %v", err)
					}

					// Add team to each repository, duplicate additions are ignored.
					addTeamToRepositories(ctx, client, filteredRepos, organization, organizationName, team, teamPermission, teamSlug)
					if err != nil {
						log.Fatalf("Failed to add team to repositories: %v", err)
					}

					return nil
				}

				// Team does not exist, create the team and add to all repositories
				filteredRepositories, err := fetchAllRepositoriesByPrefix(ctx, client, organizationName, prefix)
				if err != nil {
					log.Fatalf("Failed to fetch repositories for organization: %v", err)
				}

				// Fold repositories into repository names.
				var repoNames []string
				for _, r := range filteredRepositories {
					if strings.HasPrefix(*r.FullName, prefix) {
						repoNames = append(repoNames, *r.FullName)
					}
				}

				// More details for teams.
				description := "Read access for Git-Gopher to download logs from private repos"
				privacy := "secret"

				// Create team.
				log.Printf("Creating team %s for organization %s...", teamSlug, organizationName)
				team, r, err := client.Teams.CreateTeam(ctx.Context, organizationName, github.NewTeam{
					Name:        teamSlug,
					Description: &description,
					Permission:  &teamPermission,
					Privacy:     &privacy,
					RepoNames:   repoNames,
				})

				if r.StatusCode != 201 || err != nil {
					log.Fatalf("Could not create team for organization : %s, %v", r.Status, err)
				}

				// Add developers to team.
				for _, developer := range Developers {
					log.Printf("Adding developer %s to team %s...", developer, teamSlug)
					_, res, err := client.Teams.AddTeamMembershipByID(ctx.Context, *organization.ID, *team.ID, developer,
						&github.TeamAddTeamMembershipOptions{
							Role: "maintainer",
						})

					if res.StatusCode != 200 || err != nil {
						log.Fatalf("failed to add user to team: %v, %v", res.Status, err)
					}
				}

				addTeamToRepositories(ctx, client, filteredRepositories, organization, organizationName, team, teamPermission, teamSlug)
				if err != nil {
					log.Fatalf("Failed to add team to repositories: %v", err)
				}

				return nil
			},
		},
		{
			Name:  "analyze",
			Usage: "detect a workflow for current root",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Usage:    "path to configuation file",
					Aliases:  []string{"c"},
					Required: false,
				},
				// TODO: Complete this feature
				&cli.BoolFlag{
					Name:     "logging",
					Usage:    "enable logging, output to the set file name",
					Aliases:  []string{"l"},
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "csv",
					Usage:    "csv summary of the workflow run",
					Required: false,
				},
			},

			Subcommands: []*cli.Command{
				{
					Name:  "url",
					Usage: "analyze github project url",
					Action: func(ctx *cli.Context) error {
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

						cfg := readConfig(ctx)
						ghwf := workflow.GithubFlowWorkflow(cfg)
						violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
						if err != nil {
							log.Fatalf("Failed to analyze: %v\n", err)
						}

						workflowSummary(authors, violated, count, total, violations)

						if ctx.Bool("csv") {
							err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
							if err != nil {
								log.Fatalf("Could not create csv summary: %v", err)
							}
						}

						if ctx.Bool("logging") {
							fn, err := ghwf.WriteLog(*enrichedModel, cfg)
							if err != nil {
								log.Fatalf("Could not write json log: %v", err)
							}

							err = discord.SendLog(fn)
							if err != nil {
								log.Fatalf("Could not write json log to discord: %v", err)
							}
						}

						return nil
					},
				},
				{
					Name:    "batch",
					Aliases: []string{"b"},
					Usage:   "batch analyze a series of local git projects",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:     "csv",
							Usage:    "csv summary of the workflow run",
							Required: false,
						},
					},
					Action: func(ctx *cli.Context) error {
						utils.Environment(".env")
						path := ctx.Args().Get(0)
						if path == "" {
							path = "./"
							log.Printf("No path provided, using current directory (\"%s\") as target path", path)
						}

						var ps []string
						if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
							if info.IsDir() && info.Name() == ".git" {
								ps = append(ps, filepath.Dir(path))

								return filepath.SkipDir
							}

							return nil
						}); err != nil {
							log.Fatalf("failed to walk batch directory: %v", err)
						}
						if len(ps) == 0 {
							log.Fatalf("Could not detect any git repositories within the directiory: \"%s\"", path)
						}

						cfg := readConfig(ctx)
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

							workflowSummary(authors, v, c, t, vs)
							nameCsv := fmt.Sprintf("batch-%s.csv", filepath.Base(path))
							if ctx.Bool("csv") {
								err = ghwf.Csv(nameCsv, enrichedModel.Name, enrichedModel.URL)
								if err != nil {
									log.Fatalf("Could not create csv summary: %v", err)
								}
							}

							if ctx.Bool("logging") {
								fn, err := ghwf.WriteLog(*enrichedModel, cfg)
								if err != nil {
									log.Fatalf("Could not write json log: %v", err)
								}
								err = discord.SendLog(fn)
								if err != nil {
									log.Fatalf("Could not write json log to discord: %v", err)
								}
							}
						}

						return nil
					},
				},
			},
		},
	}

	log.Printf("BuildVersion: %s", version.BuildVersion())
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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

	var vSb strings.Builder
	for _, v := range violations {
		vSb.WriteString(v.Display(authors))
	}
	markup.Group("Violations", vSb.String())

	var sSb strings.Builder
	for _, v := range suggestions {
		sSb.WriteString(v.Display(authors))
	}
	markup.Group("Suggestions", sSb.String())

	var aSb strings.Builder
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
		aSb.WriteString(fmt.Sprintf("%s: %d\n", login, count))
	}

	aSb.WriteString(fmt.Sprintf("violated: %d\n", v))
	aSb.WriteString(fmt.Sprintf("count: %d\n", c))
	aSb.WriteString(fmt.Sprintf("total: %d\n", t))
	markup.Group("Summary", aSb.String())
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

	md.AddLine(fmt.Sprintf("Have any feedback? Feel free to submit it [here](%s)", utils.GoogleFormURL))

	return md.Render()
}

func fetchAllRepositoriesByPrefix(
	ctx *cli.Context,
	client *github.Client,
	organisationName,
	prefix string,
) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var filteredRepos []*github.Repository

	for {
		repos, res, err := client.Repositories.ListByOrg(ctx.Context, organisationName, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repositories for organization: %w", err)
		}

		for _, r := range repos {
			if strings.HasPrefix(*r.Name, prefix) {
				filteredRepos = append(filteredRepos, r)
			}
		}

		if res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage
	}

	return filteredRepos, nil
}

func addTeamToRepositories(
	ctx *cli.Context,
	client *github.Client,
	repositories []*github.Repository,
	organization *github.Organization,
	organizationName string,
	team *github.Team,
	permission string,
	teamSlug string,
) error {
	for _, r := range repositories {
		res, err := client.Teams.AddTeamRepoByID(ctx.Context, *organization.ID, *team.ID, organizationName, *r.Name,
			&github.TeamAddTeamRepoOptions{
				Permission: permission,
			})
		if res.StatusCode != 204 || err != nil {
			return fmt.Errorf("failed to add team to repository: %w", err)
		}

		log.Printf("Added team %s to repository %s", teamSlug, *r.Name)
	}

	log.Printf("Added team %s to %d repositories", teamSlug, len(repositories))

	return nil
}
