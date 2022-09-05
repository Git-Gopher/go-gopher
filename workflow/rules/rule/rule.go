package rule

import (
	"github.com/Git-Gopher/go-gopher/model/enriched"
)

type RuleCtx struct {
	Model          *enriched.EnrichedModel
	LoginWhiteList []string
}

type RuleRun func(RuleCtx) (string, *Scores)
