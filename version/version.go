package version

import (
	"fmt"
)

// Link time variables.
var (
	CommitHash  = "n/a"
	CompileDate = "n/a"
	Version     = ""
)

// BuildVersion combines available information to a nicer looking version string.
func BuildVersion() string {
	if Version != "" {
		return Version
	}

	return fmt.Sprintf("%s-%s", CommitHash, CompileDate)
}
