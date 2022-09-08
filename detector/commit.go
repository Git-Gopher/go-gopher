package detector

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
	log "github.com/sirupsen/logrus"
)

var ErrNilCommits = errors.New("nil commits")

// Write a detector/rule that can be run against commits.
type CommitDetect func(c *common, commit *local.Commit) (bool, []violation.Violation, error)

// Select a subset of the input set of commit to run the detector on.
type CommitGather func(em *enriched.EnrichedModel) ([]local.Commit, error)

// Commit detetor encapsulates the information from running a CommitDetect.
type CommitDetector struct {
	// Name of the detector.
	name string
	// Number of violations of a detector
	violated int
	// Number of units correctly following the detector.
	found int
	// Number of objects the detector acts on.
	total int
	// Particular violations that the detector reports.
	violations []violation.Violation

	gather CommitGather
	detect CommitDetect
}

func NewCommitDetector(name string, detect CommitDetect, gather CommitGather) *CommitDetector {
	return &CommitDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		gather:     gather,
		detect:     detect,
	}
}

func (cd *CommitDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return nil
	}

	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)

	c, err := NewCommon(em)
	if err != nil {
		log.Printf("could not create common: %v", err)
	}

	set, err := cd.gather(em)
	if err != nil {
		return fmt.Errorf("failed to gather commits: %w", err)
	}

	for _, co := range set {
		co := co
		detected, violations, err := cd.detect(c, &co)
		cd.total++
		if err != nil {
			return err
		}
		if detected {
			cd.found++
		}
		if violations != nil {
			cd.violations = append(cd.violations, violations...)
		}
	}

	return nil
}

func (cd *CommitDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, cd.violations
}

func (cd *CommitDetector) Name() string {
	return cd.name
}

// All commits on the main branch for github flow should be merged in,
// meaning that they have two parents(the main branch and the feature branch).
func BranchCommitDetect() (string, CommitDetect, CommitGather) {
	return "BranchCommitDetect",
		// Detector
		func(c *common, commit *local.Commit) (bool, []violation.Violation, error) {
			if len(commit.ParentHashes) >= 2 {
				return true, nil, nil
			}

			return false, nil, nil
		},
		// Working set
		func(em *enriched.EnrichedModel) ([]local.Commit, error) {
			if em.Commits == nil {
				return nil, ErrNilCommits
			}

			return em.Commits, nil
		}
}

func DiffMatchesMessageDetect() (string, CommitDetect, CommitGather) {
	return "DiffMatchesMessageDetect", func(c *common, commit *local.Commit) (bool, []violation.Violation, error) {
			words := strings.Split(commit.Message, " ")
			for _, diff := range commit.DiffToParents {
				for _, word := range words {
					all := diff.Addition + diff.Deletion + diff.Equal + diff.Name
					if strings.Contains(strings.ToLower(all), strings.ToLower(word)) {
						return true, nil, nil
					}
				}
			}

			return false, []violation.Violation{violation.NewDescriptiveCommitViolation(
				markup.Commit{
					Hash: hex.EncodeToString(commit.Hash[:]),
					GitHubLink: markup.GitHubLink{
						Owner: c.owner,
						Repo:  c.repo,
					},
				},
				commit.Message,
				commit.Committer.Email,
				commit.Committer.When,
				c.IsCurrentCommit(commit.Hash),
			)}, nil
		},
		func(em *enriched.EnrichedModel) ([]local.Commit, error) {
			if em.Commits == nil {
				return nil, ErrNilCommits
			}

			return em.Commits, nil
		}
}

// UnresolvedDetect checks if a commit is unresolved.
func UnresolvedDetect() (string, CommitDetect, CommitGather) {
	return "UnresolvedDetect", func(c *common, commit *local.Commit) (bool, []violation.Violation, error) {
			for _, diff := range commit.DiffToParents {
				lines := strings.Split(strings.ReplaceAll(diff.Addition, "\r\n", "\n"), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "<<<<<<") {
						return true, []violation.Violation{violation.NewUnresolvedMergeViolation(
							markup.Line{
								File: markup.File{
									Commit: markup.Commit{
										GitHubLink: markup.GitHubLink{
											Owner: c.owner,
											Repo:  c.repo,
										},
										Hash: hex.EncodeToString(commit.Hash.ToByte()),
									},
									Filepath: diff.Name,
								},
								Start: int(diff.Points[0].NewPosition),
								End:   nil,
							},
							commit.Committer.Email,
							commit.Committer.When,
							c.IsCurrentCommit(commit.Hash),
						)}, nil
					}
				}
			}

			return false, nil, nil
		},
		// Commit set
		func(em *enriched.EnrichedModel) ([]local.Commit, error) {
			if em.Commits == nil {
				return nil, ErrNilCommits
			}

			return em.Commits, nil
		}
}

// Check if commit is less than 3 words.
func ShortCommitMessageDetect() (string, CommitDetect, CommitGather) {
	return "ShortCommitMessageDetect", func(c *common, commit *local.Commit) (bool, []violation.Violation, error) {
			exclusions := []string{
				"first commit",
				"initial commit",
			}
			for _, exclusion := range exclusions {
				if strings.ToLower(commit.Message) == exclusion {
					return true, nil, nil
				}
			}

			if len(commit.Hash) == 0 {
				return false, nil, nil
			}

			words := strings.Split(commit.Message, " ")
			if len(words) < 5 {
				return false, []violation.Violation{violation.NewShortCommitViolation(
					markup.Commit{
						Hash: hex.EncodeToString(commit.Hash.ToByte()),
						GitHubLink: markup.GitHubLink{
							Owner: c.owner,
							Repo:  c.repo,
						},
					},
					commit.Message,
					commit.Committer.Email,
					commit.Committer.When,
					c.IsCurrentCommit(commit.Hash),
				)}, nil
			}

			return true, nil, nil
		},
		// Commit set
		func(em *enriched.EnrichedModel) ([]local.Commit, error) {
			if em.Commits == nil {
				return nil, ErrNilCommits
			}

			return em.Commits, nil
		}
}

func BinaryDetect() (string, CommitDetect, CommitGather) {
	// Extensions that should not be committed to the repository.
	disallowedExtensions := []string{".exe", ".jar", ".class"}

	return "BinaryDetect", func(c *common, commit *local.Commit) (bool, []violation.Violation, error) {
			vs := []violation.Violation{}
			for _, d := range commit.DiffToParents {
				if d.IsBinary && utils.Contains(d.Name, disallowedExtensions) {
					vs = append(vs, violation.NewBinaryViolation(
						markup.File{
							Commit: markup.Commit{
								GitHubLink: markup.GitHubLink{
									Owner: c.owner,
									Repo:  c.repo,
								},
								Hash: hex.EncodeToString(commit.Hash.ToByte()),
							},
							Filepath: d.Name,
						},
						commit.Committer.Email,
						commit.Committer.When,
						c.IsCurrentCommit(commit.Hash),
					))
				}
			}

			return len(vs) > 0, vs, nil
		},
		// Commit set
		func(em *enriched.EnrichedModel) ([]local.Commit, error) {
			if em.Commits == nil {
				return nil, ErrNilCommits
			}

			return em.Commits, nil
		}
}

func EmptyCommitDetect() (string, CommitDetect, CommitGather) {
	return "EmptyCommitDetect", func(c *common, commit *local.Commit) (bool, []violation.Violation, error) {
			isEmpty := true
			for _, d := range commit.DiffToParents {
				if d.Addition != "" || d.Deletion != "" && d.Equal != "" {
					isEmpty = false

					break
				}
			}

			vs := []violation.Violation{}
			if isEmpty {
				vs = append(vs, violation.NewEmptyCommitViolation(
					markup.Commit{
						Hash: commit.Hash.HexString(),
						GitHubLink: markup.GitHubLink{
							Owner: c.owner,
							Repo:  c.repo,
						},
					},
					commit.Committer.Email,
					commit.Committer.When,
					c.IsCurrentCommit(commit.Hash),
				))
			}

			return isEmpty, vs, nil
		},
		// Commit set
		func(em *enriched.EnrichedModel) ([]local.Commit, error) {
			if em.Commits == nil {
				return nil, ErrNilCommits
			}

			return em.Commits, nil
		}
}
