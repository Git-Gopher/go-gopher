package utils

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
	return c.Hash
}

func (c Commit) Link() string {
	return termlink.Link(c.String(), fmt.Sprintf("%s/%s", c.GitHubLink.Link(), c.Hash))
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
	return termlink.Link(b.String(), fmt.Sprintf("%s/tree/%s", b.GitHubLink.Link(), b.Name))
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
	return termlink.Link(p.String(), fmt.Sprintf("%s/pull/%s", p.GitHubLink.Link(), p.String()))
}
