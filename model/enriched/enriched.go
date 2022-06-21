package enriched

type Commit struct{}

type EnrichedModel struct {
	Commits []Commit
}

// Create an enriched model by merging the local and GitHub model.
// func NewEnrichedModel(local local.GitModel, github github.GithubModel) (*EnrichedModel, error) {
// 	return nil, nil
// }

// Create an enriched model by pulling down the repo and scraping.
// func CreateEnrichedModel(remote string) (*EnrichedModel, error) {
// 	return nil, nil
// }
