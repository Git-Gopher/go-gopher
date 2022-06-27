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
