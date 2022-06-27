package markup

import (
	"fmt"
)

// Helper functions for GitHub actions commands
func Group(title, content string) {
	fmt.Printf("::group::%s\n", title)
	fmt.Println(content)
	fmt.Println("::endgroup::")
}

func Outputs(name, value string) {
	fmt.Printf("::set-output name=%s::%s\n", name, value)
}
