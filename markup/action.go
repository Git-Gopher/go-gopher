package markup

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

const MAX_LENGTH = 65536

// Helper functions for GitHub actions commands.
func Group(title, content string) {
	fmt.Printf("::group::%s\n", title) //nolint: forbidigo
	fmt.Println(content)               //nolint: forbidigo
	fmt.Println("::endgroup::")        //nolint: forbidigo
}

func Outputs(name, value string) {
	// Sanitize input with literals
	name = strings.ReplaceAll(name, "%", `%25`)
	value = strings.ReplaceAll(value, "%", `%25`)
	name = strings.ReplaceAll(name, "\n", `%0A`)
	value = strings.ReplaceAll(value, "\n", `%0A`)
	name = strings.ReplaceAll(name, "\r", `%0D`)
	value = strings.ReplaceAll(value, "\r", `%0D`)

	// Comment max length

	if len(value) > MAX_LENGTH {
		log.Warnf("Github Action output string length is greater than %d", MAX_LENGTH)
		value = value[:MAX_LENGTH]
	}

	fmt.Printf("::set-output name=%s::%s\n", name, value) //nolint: forbidigo
}
