package markup

import (
	"testing"
)

func TestMarkdown(t *testing.T) {
	m := NewMarkdown()
	s := m.
		Title("Title").
		SubTitle("SubTitle").
		Text("Text").
		Code("Code").
		Collapsible("Collapsable", NewMarkdown().Text("Collapsible")).
		Table([]string{"A", "B", "C"}, []string{"D", "E", "F"}).
		String()
	t.Logf("\n%s", s)
}
