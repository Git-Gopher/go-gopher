# Marking

This document is intended to be used as a reference for marking tool `cmd/marker`

## Install

### Download latest build

https://github.com/Git-Gopher/go-gopher/actions/workflows/git-gopher.yml

Click into latest successful build and download artifects.

TODO: install go-releaser

### Compile locally

The tool can be compiled by using the build command
```bash
make build
```
This command will compile the tool and generate the binary `bin/marker`

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