package model

import "github.com/shurcooL/githubv4"

type Author struct {
	Login     githubv4.String
	AvatarUrl githubv4.String
}

type GithubModel struct {
	Author
}
