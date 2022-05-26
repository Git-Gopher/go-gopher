package model

import "github.com/go-git/go-git/v5/plumbing/object"

type Commit struct {
	Commit *object.Commit // remove future

	Parent  string
	Hash    string
	Message string
	Content string
}
