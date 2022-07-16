package markup

import (
	"fmt"

	"github.com/savioxavier/termlink"
)

// GitHubLink - GitHub link.
type GitHubLink struct {
	Owner string
	Repo  string
}

func (g GitHubLink) String() string {
	return fmt.Sprintf("%s/%s", g.Owner, g.Repo)
}

func (g GitHubLink) Link() string {
	return termlink.Link(g.String(), fmt.Sprintf("https://github.com/%s", g.String()))
}

// Author - Author link.
type Author string

func (a Author) String() string {
	return string(a)
}

func (a Author) Link() string {
	return termlink.Link(a.String(), fmt.Sprintf("https://github.com/%s", a.String()))
}

// Commit - Commit link.
type Commit struct {
	GitHubLink
	Hash string
}

func (c Commit) String() string {
	return c.Hash[:7] // short commit hash with 7 characters.
}

func (c Commit) Link() string {
	return termlink.Link(c.String(), fmt.Sprintf("https://github.com/%s/commit/%s", c.GitHubLink.String(), c.Hash))
}

// Branch - Branch link.
type Branch struct {
	GitHubLink
	Name string
}

func (b Branch) String() string {
	return b.Name
}

func (b Branch) Link() string {
	return termlink.Link(b.String(), fmt.Sprintf("https://github.com/%s/tree/%s", b.GitHubLink.String(), b.Name))
}

// PR - Pull request link.
type PR struct {
	GitHubLink
	Number int
}

func (p PR) String() string {
	return fmt.Sprintf("#%d", p.Number)
}

func (p PR) Link() string {
	return termlink.Link(p.String(), fmt.Sprintf("https://github.com/%s/pull/%d", p.GitHubLink.String(), p.Number))
}

// File - File link.
type File struct {
	Commit   Commit
	Filepath string
}

func (f File) String() string {
	return fmt.Sprintf("%s:%s", f.Commit.String(), f.Filepath)
}

func (f File) Link() string {
	return termlink.Link(
		f.String(),
		fmt.Sprintf("https://github.com/%s/tree/%s/%s",
			f.Commit.GitHubLink.String(),
			f.Commit.Hash,
			f.Filepath,
		),
	)
}

// Line - Line link.
type Line struct {
	File File

	Start int
	End   *int
}

func (l Line) String() string {
	if l.End == nil {
		return fmt.Sprintf(
			"%s#L%d",
			l.File.String(),
			l.Start,
		)
	}

	return fmt.Sprintf(
		"%s#L%d-L%d",
		l.File.Filepath,
		l.Start,
		l.End,
	)
}

func (l Line) Link() string {
	if l.End == nil {
		return termlink.Link(
			l.String(),
			fmt.Sprintf(
				"https://github.com/%s/tree/%s/%s#L%d",
				l.File.Commit.GitHubLink.String(),
				l.File.Commit.Hash,
				l.File.Filepath,
				l.Start,
			),
		)
	}

	return termlink.Link(
		l.String(),
		fmt.Sprintf(
			"https://github.com/%s/tree/%s/%s#L%d-L%d",
			l.File.Commit.GitHubLink.String(),
			l.File.Commit.Hash,
			l.File.Filepath,
			l.Start,
			l.End,
		),
	)
}
