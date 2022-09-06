package commands

import "github.com/urfave/cli/v2"

var (
	AnalyzeUrlCommand = &cli.Command{
		Name:      "url",
		Aliases:   []string{"u"},
		Category:  "Marker",
		Usage:     "grade a single repository with GitHub URL",
		UsageText: "go-gopher-marker url <url> - grade repository with GitHub URL",
		ArgsUsage: "<url>",
		Action:    LoadFlags(Cmd.SingleUrlCommand),
	}

	AnalyzeLocalCommand = &cli.Command{
		Name:      "local",
		Aliases:   []string{"l"},
		Category:  "Marker",
		Usage:     "grade a single local repository",
		UsageText: "go-gopher-marker local <path> - grade local repository",
		ArgsUsage: "<folder>",
		Action:    LoadFlags(Cmd.SingleLocalCommand),
	}

	AnalyzeFolderCommand = &cli.Command{
		Name:      "folder",
		Aliases:   []string{"f"},
		Category:  "Marker",
		Usage:     "grade a folder of repositories",
		UsageText: "go-gopher-marker folder <path> - grade folder of repositories",
		ArgsUsage: "<folder>",
		Action:    LoadFlags(Cmd.FolderLocalCommand),
	}
	GenerateConfigCommand = &cli.Command{
		Name:     "generate",
		Category: "Utils",
		Usage:    "generate and reset configuration files. options.yml and .env files",
		Action:   SkipFlags(Cmd.GenerateConfigCommand),
	}
)
