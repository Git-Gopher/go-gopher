package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

// var (
// 	cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
// 	memprofile = flag.String("memprofile", "mem.prof", "write memory profile to `file`")
// )

func main() {
	var data []*Sample

	file, err := os.ReadFile("./cmd/profiler/data.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	for i, repo := range data {
		data[i] = runSample(repo)

		file, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile("./cmd/profiler/result.json", file, 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runSample(sample *Sample) *Sample {
	if sample == nil {
		return nil
	}

	name := strings.Split(sample.URL, "/")[4]

	cpuprofile := fmt.Sprintf("prof/%s-cpu.prof", name)
	memprofile := fmt.Sprintf("prof/%s-mem.prof", name)

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Panicf("could not create CPU profile: %s", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Errorf("could not close CPU profile: %s", err)
			}
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Panicf("could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})

	r, err := fetch(sample.URL)
	if err != nil {
		log.Errorf("failed to fetch: %s", err)

		return sample
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Panicf("could not create memory profile: %s", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Errorf("could not close memory profile: %s", err)
			}
		}()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Panicf("could not write memory profile: %s", err)
		}
	}

	sample.CherryPick = r.CherryPick
	sample.CherryPickRelease = r.CherryPickRelease
	sample.CrissCrossMerged = r.CrissCrossMerged
	sample.FeatureBranching = r.FeatureBranching
	sample.Hotfix = r.Hotfix
	sample.Unresolved = r.Unresolved

	return sample
}

func runRepo(repo *Repository) *Repository {
	if repo == nil {
		return nil
	}

	cpuprofile := fmt.Sprintf("prof/%s-cpu.prof", repo.Name)
	memprofile := fmt.Sprintf("prof/%s-mem.prof", repo.Name)

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Panicf("could not create CPU profile: %s", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Errorf("could not close CPU profile: %s", err)
			}
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Panicf("could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})

	r, err := fetch(repo.URL)
	if err != nil {
		log.Errorf("failed to fetch: %s", err)

		return repo
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Panicf("could not create memory profile: %s", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Errorf("could not close memory profile: %s", err)
			}
		}()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Panicf("could not write memory profile: %s", err)
		}
	}

	repo.CherryPick = r.CherryPick
	repo.CherryPickRelease = r.CherryPickRelease
	repo.CrissCrossMerged = r.CrissCrossMerged
	repo.FeatureBranching = r.FeatureBranching
	repo.Hotfix = r.Hotfix
	repo.Unresolved = r.Unresolved

	return repo
}

func fetch(githubURL string) (*Repository, error) {
	utils.Environment(".env")

	// Clone repository into memory.
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: githubURL,
		Auth: &githttp.BasicAuth{
			Username: "non-empty",
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner and repo name: %w", err)
	}

	// Create enrichedModel.
	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to create enriched model: %w", err)
	}

	scoresMap := workflow.Detect(rule.RuleCtx{
		Model:          enrichedModel,
		LoginWhiteList: []string{}, // TODO: add whitelist
	})

	log.Infof("Finished running rules for %s/%s", repoOwner, repoName)

	r := &Repository{}

	for name, scores := range scoresMap {
		switch name {
		case "Cherry Pick":
			r.CherryPick = ptr(scores.GitFlow.Value())
		case "Cherry Pick Release":
			r.CherryPickRelease = ptr(scores.GitFlow.Value())
		case "Criss Cross Merged":
			r.CrissCrossMerged = ptr(scores.GitFlow.Value())
		case "Feature Branching":
			r.FeatureBranching = ptr(scores.GitFlow.Value())
		case "Hotfix":
			r.Hotfix = ptr(scores.GitFlow.Value())
		case "Unresolved":
			r.Unresolved = ptr(scores.GitFlow.Value())
		}
	}

	return r, nil
}

func ptr(n float64) *float64 {
	return &n
}
