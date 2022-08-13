# GitHub-Action

## How to install
Add a new GitHub-Action by adding a new file into `.github/git-gopher.yml`

```yaml
name: git-gopher
on:
  pull_request:
    branches: main

env:
  GITHUB_TOKEN: ${{secrets.GITHUB_TOKEN}}

jobs:
  action:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Use Cache
        uses: actions/cache@v3
        with:
          path: cache.json
          key: ${{ github.ref }}
      - name: Analyze Workflow
        id: go-gopher
        run: |
          gh release download -R https://github.com/Git-Gopher/go-gopher-action -p 'go-gopher-github-action'
          chmod +x ./go-gopher-github-action
          ./go-gopher-github-action
        env:
          GITHUB_TOKEN: ${{secrets.GITHUB_TOKEN}}
      - name: Add PR Comment
        uses: marocchino/sticky-pull-request-comment@v2
        with:
          message: "${{steps.go-gopher.outputs.pr_summary}}"
      - name: Artifact outputs
        uses: actions/upload-artifact@v3
        with:
          path: |
            log-**.json
            cache.json
```

The action is setup to only run on PRs.