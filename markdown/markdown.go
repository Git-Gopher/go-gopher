package markdown

import "strings"

// Chainable markdown generator.
type Markdown struct {
	builder *strings.Builder
}

func NewMarkdown() *Markdown {
	return &Markdown{
		builder: &strings.Builder{},
	}
}

func (m *Markdown) Title(title string) *Markdown {
	m.builder.WriteString("# " + title + "\n")

	return m
}

func (m *Markdown) SubTitle(title string) *Markdown {
	m.builder.WriteString("## " + title + "\n")

	return m
}

func (m *Markdown) Text(text string) *Markdown {
	m.builder.WriteString(text + "\n")

	return m
}

func (m *Markdown) Code(code string) *Markdown {
	m.builder.WriteString("```\n" + code + "\n```\n")

	return m
}

func (m *Markdown) Table(rows ...[]string) *Markdown {
	m.builder.WriteString("| ")
	for r, row := range rows {
		for _, cell := range row {
			m.builder.WriteString(cell + " | ")
		}
		switch r {
		case 0:
			m.builder.WriteString("\n|")
			m.builder.WriteString(strings.Repeat(" - |", len(row)))
			m.builder.WriteString("\n| ")

		case len(rows) - 1:
			m.builder.WriteString("\n")

		default:
			m.builder.WriteString("\n| ")
		}
	}

	return m
}

func (m *Markdown) Collapsible(md *Markdown) *Markdown {
	m.builder.WriteString("<details>\n")
	m.builder.WriteString("<summary>Show/Hide</summary>\n")
	m.builder.WriteString(md.String())
	m.builder.WriteString("</details>\n")

	return m
}

func (m *Markdown) String() string {
	return m.builder.String()
}
