package version

import (
	"fmt"
)

// Link time variables.
var (
	CommitHash  = ""
	CompileDate = ""
)

// BuildVersion combines available information to a nicer looking version string.
func BuildVersion() string {
	return fmt.Sprintf("%s-%s", CommitHash, CompileDate)
}
