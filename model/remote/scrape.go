package remote

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/google/go-github/v47/github"
	"github.com/shurcooL/githubv4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	ErrQueryParameters = errors.New("invalid query parameters")
	githubQuerySize    = 100
)

const GITHUB_NOREPLY_EMAIL = "noreply@github.com"

type Scraper struct {
	Client *githubv4.Client
	API    *github.Client
}

func NewScraper() Scraper {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	api := github.NewClient(httpClient)

	return Scraper{
		client,
		api,
	}
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   githubv4.String
}

func (s *Scraper) FetchPullRequestClosingIssues(
	ctx context.Context,
	owner,
	name string,
	number int,
	cursor string,
) ([]*Issue, error) {
	var q struct {
		Repository struct {
			PullRequest struct {
				ClosingIssuesReferences struct {
					Nodes []struct {
						Id     string
						Title  string
						Body   string
						Author struct {
							Login     string
							AvatarUrl string
							User      struct {
								Email string
							} `graphql:"... on User"`
						}
					}
					PageInfo PageInfo
				} `graphql:"closingIssuesReferences(first: $first, after: $cursor)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*Issue
	variables := map[string]interface{}{
		"number": githubv4.Int(number),
		"first":  githubv4.Int(githubQuerySize),
		"cursor": githubv4.String(cursor),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(ctx, q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch additional pull request closing issues references: %w", err)
		}

		for _, i := range q.Repository.PullRequest.ClosingIssuesReferences.Nodes {
			issue := Issue{
				Id:    i.Id,
				Title: i.Title,
				Body:  i.Body,
				Author: &Author{
					Login:     i.Author.Login,
					AvatarUrl: i.Author.AvatarUrl,
					Email:     i.Author.User.Email,
				},
			}

			all = append(all, &issue)
		}

		if !q.Repository.PullRequest.ClosingIssuesReferences.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.PullRequest.ClosingIssuesReferences.PageInfo.EndCursor)
	}

	return all, nil
}

func (s *Scraper) FetchIssues(ctx context.Context, owner, name string) ([]*Issue, error) {
	var q struct {
		Repository struct {
			Issues struct {
				Nodes []struct {
					Id          string
					Number      int
					Title       string
					Body        string
					State       string
					StateReason string
					// Author
					Author struct {
						Login     string
						AvatarUrl string
						User      struct {
							Email string
						} `graphql:"... on User"`
					}
				}
				PageInfo PageInfo
			} `graphql:"issues(first: $first, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*Issue
	variables := map[string]interface{}{
		"first":  githubv4.Int(githubQuerySize),
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch issues: %w", err)
		}

		for _, is := range q.Repository.Issues.Nodes {
			all = append(all, &Issue{
				Id:          is.Id,
				Number:      is.Number,
				Title:       is.Title,
				Body:        is.Body,
				State:       is.State,
				StateReason: is.StateReason,
				Author: &Author{
					Login:     is.Author.Login,
					AvatarUrl: is.Author.AvatarUrl,
					Email:     is.Author.User.Email,
				},
			})
		}

		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
	}

	return all, nil
}

// Fetch all pull requests associated with the repository.
func (s *Scraper) FetchPullRequests(ctx context.Context, owner, name string) ([]*PullRequest, error) {
	var q struct {
		Repository struct {
			PullRequests struct {
				Nodes []struct {
					// PullRequest
					Id             string
					Number         int
					HeadRefName    string
					BaseRefName    string
					Title          string
					Body           string
					ClosedAt       string
					CreatedAt      string
					ReviewDecision string
					Merged         bool
					Closed         bool
					MergedBy       struct {
						Login     string
						AvatarUrl string
						User      struct {
							Email string
						} `graphql:"... on User"`
					}
					Url string
					// Author
					Author struct {
						Login     string
						AvatarUrl string
						User      struct {
							Email string
						} `graphql:"... on User"`
					}
					// Issues
					ClosingIssuesReferences struct {
						Nodes []struct {
							Id     string
							Title  string
							Body   string
							Author struct {
								Login     string
								AvatarUrl string
								User      struct {
									Email string
								} `graphql:"... on User"`
							}
						}
						PageInfo PageInfo
					} `graphql:"closingIssuesReferences(first: $first)"`
					// XXX: Limitation of the GraphQL API, can't properly paginate reviewthreads, limiting to 100 for now.
					ReviewThreads struct {
						Nodes []struct {
							Id         string
							IsOutdated bool
							IsResolved bool
							Path       string
						}
					} `graphql:"reviewThreads(first: 100)"`
				}
				PageInfo PageInfo
			} `graphql:"pullRequests(first: $first, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*PullRequest
	variables := map[string]interface{}{
		"first":  githubv4.Int(githubQuerySize),
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch pull requests: %w", err)
		}

		for _, mpr := range q.Repository.PullRequests.Nodes {
			// ISO8061 layout.
			layout := "2006-01-02T15:04:05Z0700"

			var createdAt, closedAt *time.Time

			// Pull request will always have valid creation date.
			{
				created, err := time.Parse(layout, mpr.CreatedAt)
				if err != nil {
					return nil, fmt.Errorf("could not parse ISO time for PR creation: %w", err)
				}
				createdAt = &created
			}

			// Pull request has not been closed.
			if mpr.ClosedAt == "" {
				closedAt = nil
			} else {
				closed, err := time.Parse(layout, mpr.ClosedAt)
				if err != nil {
					return nil, fmt.Errorf("could not parse ISO time for PR close: %w", err)
				}
				closedAt = &closed
			}

			pr := PullRequest{
				Id:             mpr.Id,
				Number:         mpr.Number,
				HeadRefName:    mpr.HeadRefName,
				BaseRefName:    mpr.BaseRefName,
				CreatedAt:      createdAt,
				ClosedAt:       closedAt,
				Title:          mpr.Title,
				Body:           mpr.Body,
				Closed:         mpr.Closed,
				ReviewDecision: mpr.ReviewDecision,
				Merged:         mpr.Merged,
				MergedBy: &Author{
					Login:     mpr.MergedBy.Login,
					AvatarUrl: mpr.MergedBy.AvatarUrl,
					Email:     mpr.MergedBy.User.Email,
				},
				Url: mpr.Url,
				Author: &Author{
					Login:     mpr.Author.Login,
					AvatarUrl: mpr.Author.AvatarUrl,
					Email:     mpr.Author.User.Email,
				},
				ClosingIssues: nil,
				ReviewThreads: nil,
			}

			// Closing issues
			var cis []*Issue = make([]*Issue, len(mpr.ClosingIssuesReferences.Nodes))
			for i, ci := range mpr.ClosingIssuesReferences.Nodes {
				cis[i] = &Issue{
					Id:    ci.Id,
					Title: ci.Title,
					Body:  ci.Body,
					Author: &Author{
						Login:     ci.Author.Login,
						AvatarUrl: ci.Author.AvatarUrl,
						Email:     ci.Author.User.Email,
					},
				}
			}

			if mpr.ClosingIssuesReferences.PageInfo.HasNextPage {
				acir, err := s.FetchPullRequestClosingIssues(ctx, owner, name, pr.Number,
					string(mpr.ClosingIssuesReferences.PageInfo.EndCursor))
				if err != nil {
					return nil, fmt.Errorf("Failed to fetch pull request closing issue references: %w", err)
				}

				cis = append(cis, acir...)
			}

			pr.ClosingIssues = cis

			// Review threads
			var rs []*ReviewThread = make([]*ReviewThread, len(mpr.ReviewThreads.Nodes))
			for i, r := range mpr.ReviewThreads.Nodes {
				rt := ReviewThread{
					Id:         r.Id,
					IsResolved: r.IsResolved,
					IsOutdated: r.IsOutdated,
					Path:       r.Path,
				}

				rs[i] = &rt
			}

			pr.ReviewThreads = rs

			all = append(all, &pr)
		}

		if !q.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.PullRequests.PageInfo.EndCursor)
	}

	return all, nil
}

// Fetch basic information that doesn't require pagnation (url).
// Although this can be constructed from the owner and name I would rather fetch it.
func (s *Scraper) FetchURL(ctx context.Context, owner, name string) (string, error) {
	var q struct {
		Repository struct {
			Url string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	if err := s.Client.Query(ctx, &q, variables); err != nil {
		return "", fmt.Errorf("Failed to fetch pull requests: %w", err)
	}

	return q.Repository.Url, nil
}

// FetchCommitters, get all committers from a repo.
func (s *Scraper) FetchCommitters(ctx context.Context, owner, name string) ([]Committer, error) {
	committers, err := s.FetchCommittersDefaultBranch(ctx, owner, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch committers from default branch: %w", err)
	}

	headName, branchNames, err := s.FetchBranchHeads(ctx, owner, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch branch heads: %w", err)
	}

	errCh := make(chan error)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for _, branchName := range branchNames {
		wg.Add(2)
		go func(branchName string) {
			wg.Done()
			compare, _, err := s.API.Repositories.CompareCommits(ctx, owner, name, headName, branchName, nil)
			if err != nil {
				errCh <- fmt.Errorf("Failed to compare commits: %w", err)

				return
			}

			branchCommitters, err := s.FetchCommittersBranch(ctx, owner, name, branchName, *compare.AheadBy)
			if err != nil {
				errCh <- fmt.Errorf("Failed to fetch committers from branch: %w", err)

				return
			}

			mutex.Lock()
			committers = append(committers, branchCommitters...)
			mutex.Unlock()

			wg.Done()
		}(branchName)
	}

	select {
	case err := <-errCh:
		return nil, err
	default:
		wg.Wait()
	}

	return committers, nil
}

// FetchCommittersDefaultBranch, get all committers from default branch of a repo.
func (s *Scraper) FetchCommittersDefaultBranch(ctx context.Context, owner, name string) ([]Committer, error) {
	var q struct {
		Repository struct {
			DefaultBranchRef struct {
				Target struct {
					Commit struct {
						History struct {
							Nodes []struct {
								Oid    string
								Author struct {
									Email string
									User  struct {
										Login string
									}
								}
								Committer struct {
									Email string
									User  struct {
										Login string
									}
								}
							}
							PageInfo PageInfo
						} `graphql:"history(first: $first, after: $cursor)"`
					} `graphql:"... on Commit"`
				} `graphq:"target"`
			} `graphql:"defaultBranchRef"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []Committer
	variables := map[string]interface{}{
		"first":  githubv4.Int(githubQuerySize),
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch committers: %w", err)
		}

		for _, i := range q.Repository.DefaultBranchRef.Target.Commit.History.Nodes {
			if i.Author.Email != GITHUB_NOREPLY_EMAIL {
				committer := Committer{
					CommitId: i.Oid,
					Email:    i.Author.Email,
					Login:    i.Author.User.Login,
				}
				all = append(all, committer)
			}

			if i.Author.Email != i.Committer.Email && i.Committer.Email != GITHUB_NOREPLY_EMAIL {
				committer := Committer{
					CommitId: i.Oid,
					Email:    i.Committer.Email,
					Login:    i.Committer.User.Login,
				}
				all = append(all, committer)
			}
		}

		if !q.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
	}

	return all, nil
}

// FetchCommittersDefaultBranch, get all committers from default branch of a repo.
func (s *Scraper) FetchCommittersBranch(
	ctx context.Context,
	owner, name, branch string,
	limit int,
) ([]Committer, error) {
	if limit > 100 {
		limit = 100
	}
	var q struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commit struct {
						History struct {
							Nodes []struct {
								Oid    string
								Author struct {
									Email string
									User  struct {
										Login string
									}
								}
								Committer struct {
									Email string
									User  struct {
										Login string
									}
								}
							}
							PageInfo PageInfo
						} `graphql:"history(first: $first)"`
					} `graphql:"... on Commit"`
				} `graphq:"target"`
			} `graphql:"ref(qualifiedName: $qualifiedName)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []Committer
	variables := map[string]interface{}{
		"first":         githubv4.Int(limit),
		"owner":         githubv4.String(owner),
		"name":          githubv4.String(name),
		"qualifiedName": githubv4.String(branch),
	}

	if err := s.Client.Query(ctx, &q, variables); err != nil {
		return nil, fmt.Errorf("Failed to fetch committers: %w", err)
	}

	for _, i := range q.Repository.Ref.Target.Commit.History.Nodes {
		if i.Author.Email != GITHUB_NOREPLY_EMAIL {
			committer := Committer{
				CommitId: i.Oid,
				Email:    i.Author.Email,
				Login:    i.Author.User.Login,
			}
			all = append(all, committer)
		}

		if i.Author.Email != i.Committer.Email && i.Committer.Email != GITHUB_NOREPLY_EMAIL {
			committer := Committer{
				CommitId: i.Oid,
				Email:    i.Committer.Email,
				Login:    i.Committer.User.Login,
			}
			all = append(all, committer)
		}
	}

	return all, nil
}

// FetchBranchHeads, get all heads from a repo without main branch.
func (s *Scraper) FetchBranchHeads(ctx context.Context, owner, name string) (string, []string, error) {
	var q struct {
		Repository struct {
			DefaultBranchRef struct {
				Name string
			} `graphql:"defaultBranchRef"`
			Refs struct {
				Nodes []struct {
					Name   string
					Target struct {
						Commit struct {
							PushedDate             time.Time
							AssociatedPullRequests struct {
								Nodes []struct {
									Merged bool
								}
							} `graphql:"associatedPullRequests(first: 100)"`
						} `graphql:"... on Commit"`
					} `graphq:"target"`
				} `graphql:"nodes"`
				PageInfo PageInfo
			} `graphql:"refs(first: $first, after: $cursor, refPrefix: $refPrefix)"` // does not support cursor
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"first":     githubv4.Int(githubQuerySize),
		"cursor":    (*githubv4.String)(nil),
		"owner":     githubv4.String(owner),
		"name":      githubv4.String(name),
		"refPrefix": githubv4.String("refs/heads/"),
	}

	m := []string{}

	for {
		if err := s.Client.Query(ctx, &q, variables); err != nil {
			return "", nil, fmt.Errorf("Failed to fetch branch heads: %w", err)
		}

		for _, i := range q.Repository.Refs.Nodes {
			// skip main branch
			if i.Name == q.Repository.DefaultBranchRef.Name {
				continue
			}

			// skip stale branches
			if i.Target.Commit.PushedDate.Before(time.Now().Add(-time.Hour * 24 * 30)) {
				continue
			}

			if len(i.Target.Commit.AssociatedPullRequests.Nodes) == 0 {
				m = append(m, i.Name)

				continue
			}

			for _, j := range i.Target.Commit.AssociatedPullRequests.Nodes {
				if j.Merged {
					break
				}

				m = append(m, i.Name)
			}
		}

		if !q.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(q.Repository.Refs.PageInfo.EndCursor)
	}

	return q.Repository.DefaultBranchRef.Name, m, nil
}

// Fetch repositories that have a minimum amoutn of stars, issues and contributors.
func (s *Scraper) FetchPopularRepositories( // nolint: gocognit // needs to be complex
	ctx context.Context,
	minStars int,
	minIssues int,
	minContributors int,
	minLanguages int, // helpful to filter activist or wiki type repositories
	minPullRequests int,
	numRepos int,
) ([]Repository, error) {
	if minStars < 0 || numRepos < 0 {
		return nil, fmt.Errorf("%w: %v, %v", ErrQueryParameters, minStars, numRepos)
	}

	var q struct {
		Search struct {
			Nodes []struct {
				Repository struct {
					Name       string
					Url        string
					Stargazers struct {
						TotalCount int
					}
					Languages struct {
						Nodes []struct {
							Name string
						}
					} `graphql:"languages(first: 100)"` // hopefully 100 will be enough to cover them all
					Issues struct {
						TotalCount int
					}
					PullRequests struct {
						TotalCount int
					}
					Releases struct {
						TotalCount int
					} `graphql:"Releases: refs(first: 0, refPrefix: \"refs/tags/\")"`
				} `graphql:"... on Repository"`
			}
			PageInfo PageInfo
		} `graphql:"search(first: $first, after: $cursor, query: $searchQuery, type: $searchType)"`
	}

	first := githubQuerySize
	if numRepos < first {
		first = numRepos
	}

	variables := map[string]interface{}{
		"first":  githubv4.Int(first),
		"cursor": (*githubv4.String)(nil),
		// Stars are limited between starrs..stars+10
		"searchQuery": githubv4.String(fmt.Sprintf("stars:%d..%d is:public archived:false mirror:false", minStars, minStars+40)),
		"searchType":  githubv4.SearchTypeRepository,
	}

	skippedRepos := 0
	var acceptedRepos []Repository
	for len(acceptedRepos) != numRepos {
		if err := s.Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("failed to fetch popular repositories: %w", err)
		}

		// Some query graphql search parameters simply don't exist.
		// Have to manually check if the repository is a hit or miss for a number of fields.
		// Eg: Contributors.
		for _, r := range q.Search.Nodes {
			candidateRepo := Repository{
				Name:         r.Repository.Name,
				Url:          r.Repository.Url,
				Stargazers:   r.Repository.Stargazers.TotalCount,
				Issues:       r.Repository.Issues.TotalCount,
				PullRequests: r.Repository.PullRequests.TotalCount,
			}

			languages := []string{}
			for _, l := range r.Repository.Languages.Nodes {
				languages = append(languages, l.Name)
			}

			candidateRepo.Languages = languages

			owner, name, err := utils.OwnerNameFromUrl(candidateRepo.Url)
			// Log and ignore if unable to parse
			if err != nil {
				log.Warnf("unable to parse url (%s) for owner and name %v", candidateRepo.Url, err)

				continue
			}

			// Count number of contributors via pages
			// nolint: lll
			// https://stackoverflow.com/questions/44347339/github-api-how-efficiently-get-the-total-contributors-amount-per-repository
			_, res, err := s.API.Repositories.ListContributors(ctx, owner, name,
				&github.ListContributorsOptions{
					Anon: "true",
					ListOptions: github.ListOptions{
						PerPage: 1,
					},
				},
			)
			if err != nil || res.StatusCode != 200 {
				log.Infof("unable to fetch contributors for %s, skipping...", candidateRepo.Url)

				continue
			}

			candidateRepo.Contributors = res.LastPage - res.FirstPage + 1

			if candidateRepo.Contributors >= minContributors &&
				candidateRepo.Contributors <= minContributors+10 &&
				candidateRepo.Issues >= minIssues &&
				//candidateRepo.Issues <= minIssues+50 &&
				len(candidateRepo.Languages) >= minLanguages &&
				//len(candidateRepo.Languages) <= minLanguages+4 &&
				candidateRepo.PullRequests >= minPullRequests {
				//candidateRepo.PullRequests <= minPullRequests+50 {
				acceptedRepos = append(acceptedRepos, candidateRepo)
				log.Infof("\033[32m ADDED \033[0m %s,\t %d stars, %d contributors, %d issues, %d languages, %d prs",
					candidateRepo.Url,
					candidateRepo.Stargazers,
					candidateRepo.Contributors,
					candidateRepo.Issues,
					len(candidateRepo.Languages),
					candidateRepo.PullRequests,
				)
			} else {
				log.Infof("\033[31m SKIPPING \033[0m %s,\t %d stars, %d contributors, %d issues, %d languages, %d prs",
					candidateRepo.Url,
					candidateRepo.Stargazers,
					candidateRepo.Contributors,
					candidateRepo.Issues,
					len(candidateRepo.Languages),
					candidateRepo.PullRequests,
				)
				skippedRepos++
			}
		}

		log.Printf("\033[34m collected %d repos so far, skipped %d repos (%d remaining to collect)... \033[0m", len(acceptedRepos), skippedRepos, numRepos-len(acceptedRepos))

		if !q.Search.PageInfo.HasNextPage {
			if len(acceptedRepos) < numRepos {
				log.Print("ran out of repos while scraping (out of pages)")
			} else {
				log.Print("successfully scraped all repos")
			}
			break
		}

		first := githubQuerySize
		if (numRepos - len(acceptedRepos)) < first {
			first = numRepos - len(acceptedRepos)
		}

		variables["first"] = githubv4.NewInt(githubv4.Int(first))
		variables["cursor"] = githubv4.NewString(q.Search.PageInfo.EndCursor)
	}

	return acceptedRepos, nil
}
