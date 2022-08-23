# Marking

This document is intended to be used as a reference for marking tool `cmd/marker`

## Install

### Download latest build

Releases can be accessed here https://github.com/Git-Gopher/go-gopher/releases

Latest builds can be found here https://github.com/Git-Gopher/go-gopher-action/releases/tag/latest

Click into latest successful build and download artifects.

### Compile locally

The tool can be compiled by using the build command
```bash
make build # Auto detect windows or linux or macos
```

This command will compile the tool and generate the binary 
- `bin/go-gopher-cli`
- `bin/go-gopher-marker`
- `bin/go-gopher-github`

Build for other platforms can be done by using the build command
```bash
make release-windows
make release-linux
make release-macos
```

### Run locally

Run the tool locally using Go.

Requirement
- install Go 1.18
- gcc or any C compiler installed

```bash
# Get into this repo directory
cd go-gopher 

# Install dependencies (first install)
go get
# or `go mod tidy`

# Run the tool
go run ./cmd/marker <command> <arg>

# Examples
go run ./cmd/marker help
go run ./cmd/marker url https://github.com/Git-Gopher/go-gopher
```

## Commands

| Command | Args  | Usage  |
|---|---|---|
| url  | \<url>  | run marker on a remote GitHub repo using url  |
| local   | \<folder>  | run marker on a cloned repo locally  |
| folder  | \<folder>  | run marker on multiple cloned repo |
| generate  |   | generates config files  |
| help |   | display all commands |

Examples:
```
go-gopher-marker url https://github.com/Git-Gopher/go-gopher

go-gopher-marker local ./my/git/repo

go-gopher-marker folder ./my/gits

go-gopher-marker generate

go-gopher-marker help
```

## Configurations

To generate/reset default configurations run
```
go-gopher-marker generate
```

This will generate two files for the marker: `.env` and `options.yml`

.env file should contain
```
GITHUB_TOKEN=
```
This token is required for private repos ignore if repos are public.

You can override `GITHUB_TOKEN` by providing it as a `--token` flag with cli.
```
go-gopher-marker --token <GITHUB_TOKEN> local ./my/git/repo
```

## Output of marker

The marker would generate a markdown file per student it is marking.

e.g.


TODO: set output folder
TODO: customise output template
TODO: CSV with all student marks