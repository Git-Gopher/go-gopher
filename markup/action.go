package markup

import (
	"fmt"
	"strings"
)

// Helper functions for GitHub actions commands
func Group(title, content string) {
	fmt.Printf("::group::%s\n", title)
	fmt.Println(content)
	fmt.Println("::endgroup::")
}

func Outputs(name, value string) {
	// Sanitize with literals
	name = strings.ReplaceAll(name, "\n", `\n`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	name = strings.ReplaceAll(name, "\r", `\r`)
	value = strings.ReplaceAll(value, "\r", `\r`)

	fmt.Printf("::set-output name=%s::%s\n", name, value)
}
