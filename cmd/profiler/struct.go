package main

type Repository struct {
	// data.json
	Name                 string   `json:"name"`
	URL                  string   `json:"url"`
	Stargazers           int      `json:"stargazers"`
	Languages            []string `json:"languages"`
	Issues               int      `json:"issues"`
	PullRequests         int      `json:"pullRequests"`
	Contributors         int      `json:"contributors"`
	PrimaryBranchCommits int      `json:"primaryBranchCommits"`

	// data
	CherryPick        *float64 `json:"ruleCherryPick"`
	CherryPickRelease *float64 `json:"ruleCherryPickRelease"`
	CrissCrossMerged  *float64 `json:"ruleCrissCrossMerged"`
	FeatureBranching  *float64 `json:"ruleFeatureBranching"`
	Hotfix            *float64 `json:"rulehotfix"`
	Unresolved        *float64 `json:"ruleUnresolved"`
}

type Sample struct {
	URL      string `json:"url"`
	Workflow string `json:"workflow"`

	// data
	CherryPick        *float64 `json:"ruleCherryPick"`
	CherryPickRelease *float64 `json:"ruleCherryPickRelease"`
	CrissCrossMerged  *float64 `json:"ruleCrissCrossMerged"`
	FeatureBranching  *float64 `json:"ruleFeatureBranching"`
	Hotfix            *float64 `json:"rulehotfix"`
	Unresolved        *float64 `json:"ruleUnresolved"`
}
