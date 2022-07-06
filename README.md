# go-gopher

go-gopher is a git workflow analysis tool to check cohesion to popular git workflows

## Running

1. Ensure that the project [environment](#environment) file is filled out
2. Run `make` in the root of the project directory to create the `go-gopher` binary
3. Run the binary `./go-gopher --help` to display help.

A good starting point is to navigate to a git project and run `./go-gopher analyze local ./ --csv`

## Environment

Environment variables are defined in the `.env` which should be created from the [`.env.example`](./.env.example) template.

| Environment variables | Description                                                                           |
| --------------------- | ------------------------------------------------------------------------------------- |
| GITHUB_TOKEN          | Used in accessing GitHub's GraphQL API, requires read access to repositories          |
| CACHE_PATH            | Base path to save the cache of previous runs to, as a fallback from $GITHUB_WORKSPACE |

## Configuration

Detectors can be toggled and weighted within [config/config.yml](./config/config.yml)
## GitHub Action

This project is available as a GitHub action. An example usage can be found within this [project](./.github/workflows/publish.yml). The action has been published via the [go-gopher-action](https://github.com/Git-Gopher/go-gopher-action) repository.
