package workflow

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model"
)

func TestWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		workflow *Workflow
		model    *model.GitModel
		wantErr  bool
		want     bool
	}{
		{"Github Flow", &GithubFlowWorkflow, nil, false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a model

			// Analyze the model
			tc.workflow.Analyze(nil)

			// Expected results
		})
	}
}
