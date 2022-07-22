package markup

import (
	"fmt"
	"strings"
)

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

func (m *Markdown) SubSubTitle(title string) *Markdown {
	m.builder.WriteString("### " + title + "\n")

	return m
}

func (m *Markdown) Text(text string) *Markdown {

	// Fix termlink issues
	text = strings.Replace(text, "\u001b[m", "", -1)

	m.builder.WriteString(text + "\n")

	return m
}

func (m *Markdown) Code(code string) *Markdown {
	m.builder.WriteString("```\n" + code + "\n```\n")

	return m
}

func (m *Markdown) Table(rows ...[]string) *Markdown {
	if len(rows) == 0 {
		return m
	}

	m.builder.WriteString(strings.Join(rows[0], " | "))
	m.builder.WriteString("\n")
	s := strings.Split(strings.Repeat("--- ", len(rows[0])), " ")
	m.builder.WriteString(strings.Join(s, "|"))
	m.builder.WriteString("\n")

	if len(rows) == 1 {
		s := strings.Split(strings.Repeat(" ", len(rows[0])), "")
		m.builder.WriteString(strings.Join(s, " | "))
		m.builder.WriteString("\n")

		return m
	}

	for _, row := range rows[1:] {
		m.builder.WriteString(strings.Join(row, " | "))
		m.builder.WriteString("\n")
	}

	return m
}

func (m *Markdown) Collapsible(title string, body *Markdown) *Markdown {
	m.builder.WriteString("<details>\n")
	m.builder.WriteString(fmt.Sprintf("<summary>%s</summary>\n", title))
	m.builder.WriteString(body.String())
	m.builder.WriteString("</details>\n")

	return m
}

func (m *Markdown) String() string {
	return m.builder.String()
}
