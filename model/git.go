package model

import "github.com/go-git/go-git/v5/plumbing/object"

type Commit struct {
	Commit *object.Commit // remove future

	Parent  string
	Hash    string
	Message string
	Content string
}

type Branch struct {
}

// TODO: fill this out
type GitModel struct {
	Commits []Commit
}

// At the end of the day I should be able to pass in the model from go-git (or even a url) and create our modifiied git model
func NewGitModel() *GitModel {

	return nil
}
