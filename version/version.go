package version

import (
	"fmt"
)

// Link time variables.
var (
	CommitHash  = "n/a"
	CompileDate = "n/a"
)

// BuildVersion combines available information to a nicer looking version string.
func BuildVersion() string {
	return fmt.Sprintf("%s-%s", CommitHash, CompileDate)
}
