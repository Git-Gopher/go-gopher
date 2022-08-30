# go-gopher

go-gopher is a git workflow analysis tool to check cohesion to popular git workflows

## Running

1. Ensure that the project [environment](#environment) file is filled out
2. Run `make` in the root of the project directory to create the project binaries that will output to `bin/`
3. Run the binary `./bin/go-gopher help` to display help.

A good starting point is to navigate to a git project and run the go-gopher tool: `go-gopher analyze local ./ --csv`

## Environment

Environment variables are defined in the `.env` which should be created from the [`.env.example`](.env.example) template.

| Environment variables | Description                                                                  |
| --------------------- | ---------------------------------------------------------------------------- |
| GITHUB_TOKEN          | Used in accessing GitHub's GraphQL API, requires read access to repositories |

## Configuration

Detectors can be toggled and weighted within [config/config.yml](./config/config.json)

## GitHub Action

This project is available as a GitHub action. An example usage can be found within this [project](./.github/workflows/git-gopher.yml). The action has been published via the [go-gopher-action](https://github.com/Git-Gopher/go-gopher-action) repository.

## Documentation

Documentation for best practices and development decisions can be found in [docs](./docs/)
