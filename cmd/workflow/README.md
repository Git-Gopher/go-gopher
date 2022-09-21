# workflow
In order to evaluate the workflows of selected repositories

1. Run the go-gopher-cli query with whatever requirement parameters you want

```bash
./go-gopher-cli query --json <repos.json> <numStars> <numIssues> <numContributors> <numLanguages> <numRepos>
```

Advice for `numLanguages` is that you should set this to a value >= `1` if you're working with repositories that are relatively popular. This is because popular repositories are often non-code (eg: wikis, social activism) and having a language requirement mostly filters these out.

For our research we have decided that we should be using the following params

```bash
./go-gopher-cli query --json output.json 1000 100 8 1 1000
```
ie: numStars=1000 numIssues=100 numContributors=8 numLanguages=1 numRepos=1000

2. Use json output containing repos with go-gopher-workflow

```bash
./go-gopher-workflow batch --output results.json <repos.json>
```

This clone all those repositories in memory (in parallel) and output the results to `results.json`. From there you can import the results into a notebook to process.
