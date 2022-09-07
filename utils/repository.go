package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/urfave/cli/v2"
)

// Fetch all repositories and filter them by a prefix.
func FetchAllRepositoriesByPrefix(
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

// Add a team to a repository.
func AddTeamToRepositories(
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
