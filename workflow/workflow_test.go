package workflow

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model/local"
)

// TODO: Complete this test with all workflows.
func TestWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		workflow *Workflow
		model    *local.GitModel
		wantErr  bool
		want     bool
	}{
		{"Github Flow", GithubFlowWorkflow(nil), nil, false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a model

			// Analyze the model
			// tc.workflow.Analyze(nil)

			// Expected results
		})
	}
}
