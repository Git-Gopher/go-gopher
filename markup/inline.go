package markup

import "fmt"

func Italic(text string) string {
	return fmt.Sprintf("*%s*", text)
}

func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func Strike(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

func Link(title string, link string) string {
	return fmt.Sprintf("[%s](%s)", title, link)
}

func InlineCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}
