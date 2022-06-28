package local

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestFetchBranchGraph(t *testing.T) {
	utils.Environment("../../.env")
	r := utils.FetchRepository(t, "https://github.com/Git-Gopher/tests", "test/two-parents-merged/0")

	bIter, err := r.Branches()
	if err != nil {
		t.Errorf("Branches() error = %v", err)
	}

	reference, err := bIter.Next()
	if err != nil {
		t.Errorf("bIter.Next() error = %v", err)
	}

	parent, err := r.CommitObject(reference.Hash())
	if err != nil {
		t.Errorf("CommitObject() = %v", err)
	}

	graph := FetchBranchGraph(parent)

	t.Logf("parent: %+v", parent)
	t.Logf("graph: %+v", graph)

	if graph == nil {
		t.Errorf("FetchBranchGraph() = %v", graph)
	}

	if graph.Head == nil || graph.Head.Hash == "" {
		t.Errorf("FetchBranchGraph().Head = %v", graph.Head)
	}

	expectedHash := "9513562a03d8a827a626a5a0b8cae52d48d25128"
	if graph.Head.Hash != expectedHash {
		t.Errorf("FetchBranchGraph().Head.Hash = %v, expected = %v", graph.Head.Hash, expectedHash)
	}

	if graph.Head.ParentCommits == nil || len(graph.Head.ParentCommits) != 2 {
		t.Errorf("len(FetchBranchGraph().Head.ParentCommits) = %v, expected = 2", graph.Head)
	}
}
