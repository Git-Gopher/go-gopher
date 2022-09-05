package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/google/go-github/v45/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var (
	// Developers to add to the team.
	developers = []string{"wqsz7xn", "scorpionknifes"}

	TeamComand = &cli.Command{
		Name:  "team",
		Usage: "Add developers as a team to repositories within an organization",
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
			// Setup command configuration.
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

			// Create authenticated client.
			ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			tc := oauth2.NewClient(ctx.Context, ts)
			client := github.NewClient(tc)

			log.Printf("Fetching organization %s...\n", organizationName)
			organization, _, err := client.Organizations.Get(ctx.Context, organizationName)
			if err != nil {
				log.Fatalf("Could not fetch orginization: %v", err)
			}

			// Team details.
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
			for _, developer := range developers {
				log.Printf("Adding developer %s to team %s...", developer, teamSlug)
				_, res, err := client.Teams.AddTeamMembershipByID(ctx.Context, *organization.ID, *team.ID, developer,
					&github.TeamAddTeamMembershipOptions{
						Role: "maintainer",
					})

				if res.StatusCode != 200 || err != nil {
					log.Fatalf("failed to add user to team: %v, %v", res.Status, err)
				}
			}

			// Finally add that team to the repositories.
			addTeamToRepositories(ctx, client, filteredRepositories, organization, organizationName, team, teamPermission, teamSlug)
			if err != nil {
				log.Fatalf("Failed to add team to repositories: %v", err)
			}

			return nil
		},
	}
)

// Fetch and filter all repositories belonging to an organization by a prefix string of the repository name.
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

// Add a prexisting team to a reposiroty.
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
