package markup

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
	return fmt.Sprintf("https://github.com/%s", g.String())
}

func (g GitHubLink) Markdown() string {
	return fmt.Sprintf("[%s](%s)", g.String(), g.Link())
}

// Author - Author link.
type Author string

func (a Author) String() string {
	return string(a)
}

func (a Author) Link() string {
	return fmt.Sprintf("https://github.com/%s", a.String())
}

func (a Author) Markdown() string {
	return fmt.Sprintf("[@%s](%s)", a.String(), a.Link())
}

// Commit - Commit link.
type Commit struct {
	GitHubLink
	Hash string
}

func (c Commit) String() string {
	// XXX: Quick hack, seems that the hash is missing from some structs.
	if c.Hash == "" {
		log.Panicln("[WARN] Commit hash is empty, SHOULD NEVER HAPPEN")
	}

	return c.Hash[:7] // short commit hash with 7 characters.
}

func (c Commit) Link() string {
	return fmt.Sprintf("https://github.com/%s/commit/%s", c.GitHubLink.String(), c.Hash)
}

func (c Commit) Markdown() string {
	return fmt.Sprintf("[%s](%s)", c.String(), c.Link())
}

// Branch - Branch link.
type Branch struct {
	GitHubLink
	Name string
}

func (b Branch) String() string {
	// branch name is a hash cause it has been deleted
	if len(b.Name) == 40 {
		return b.Name[:7] // short commit hash with 7 characters.
	}

	return b.Name
}

func (b Branch) Link() string {
	return fmt.Sprintf("https://github.com/%s/tree/%s", b.GitHubLink.String(), b.Name)
}

func (b Branch) Markdown() string {
	return fmt.Sprintf("[%s](%s)", b.String(), b.Link())
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
	return fmt.Sprintf("https://github.com/%s/pull/%d",
		p.GitHubLink.String(),
		p.Number)
}

func (p PR) Markdown() string {
	return fmt.Sprintf("[%s](%s)", p.String(), p.Link())
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
	return fmt.Sprintf("https://github.com/%s/tree/%s/%s", f.Commit.GitHubLink.String(), f.Commit.Hash, f.Filepath)
}

func (f File) Markdown() string {
	return fmt.Sprintf("[%s](%s)", f.String(), f.Link())
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
		return fmt.Sprintf(
			"https://github.com/%s/tree/%s/%s#L%d",
			l.File.Commit.GitHubLink.String(),
			l.File.Commit.Hash,
			l.File.Filepath,
			l.Start,
		)
	}

	return fmt.Sprintf(
		"https://github.com/%s/tree/%s/%s#L%d-L%d",
		l.File.Commit.GitHubLink.String(),
		l.File.Commit.Hash,
		l.File.Filepath,
		l.Start,
		l.End,
	)
}

func (l Line) Markdown() string {
	return fmt.Sprintf("[%s](%s)", l.String(), l.Link())
}
