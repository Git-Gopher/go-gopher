package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

const repoName = "assignment-3-and-project-team-05"

// var (
// 	cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
// 	memprofile = flag.String("memprofile", "mem.prof", "write memory profile to `file`")
// )

var (
	cpuprofile = flag.String("cpuprofile", fmt.Sprintf("prof/%s-cpu.prof", repoName), "write cpu profile to `file`")
	memprofile = flag.String("memprofile", fmt.Sprintf("prof/%s-mem.prof", repoName), "write memory profile to `file`")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
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

	err := fetch("https://github.com/GitWorkflowPractice/" + repoName)
	if err != nil {
		log.Panicf("failed to fetch: %s", err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
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
}

func fetch(githubURL string) error {
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
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return fmt.Errorf("failed to get owner and repo name: %w", err)
	}

	// Create enrichedModel.
	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("failed to create enriched model: %w", err)
	}

	// Populate authors from enrichedModel.
	authors := enriched.PopulateAuthors(enrichedModel)

	// Read marker configs
	o := LoadOptions(log.StandardLogger())
	analyzers := assess.LoadAnalyzer(o)

	start := time.Now()

	candidates := assess.RunMarker(
		analysis.MarkerCtx{
			Model:        enrichedModel,
			Contribution: analysis.NewContribution(*enrichedModel),
			Author:       authors,
		},
		analyzers,
	)

	elapsed := time.Since(start)
	log.Infof("Ran marker in %s", elapsed)

	for _, candidate := range candidates {
		log.Printf("#### @%s ####\n", candidate.Username)
	}

	return nil
}
