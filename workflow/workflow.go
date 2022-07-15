package workflow

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/config"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
)

// XXX: This is a hack to get the name of the detector.
// This should really be done using reflect so that you don't
// have to think about changing this part of the code whenever you
// change the name of a detector, making this rather brittle.
var (
	DefaultCsvPath   = "summary.csv"
	detectorRegistry = map[string]detector.Detector{
		"StaleBranchDetect":               detector.NewBranchDetector(detector.StaleBranchDetect()),
		"PullRequestApprovalDetector":     detector.NewPullRequestDetector(detector.PullRequestApprovalDetector()),
		"PullRequestIssueDetector":        detector.NewPullRequestDetector(detector.PullRequestIssueDetector()),
		"PullRequestReviewThreadDetector": detector.NewPullRequestDetector(detector.PullRequestReviewThreadDetector()),
		"DiffMatchesMessageDetect":        detector.NewCommitDetector(detector.DiffMatchesMessageDetect()),
		"ShortCommitMessageDetect":        detector.NewCommitDetector(detector.ShortCommitMessageDetect()),
		"DiffDistanceCalculation":         detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation()),
		"NewBranchNameConsistencyDetect":  detector.NewBranchCompareDetector(detector.NewBranchNameConsistencyDetect()),
		"NewFeatureBranchDetector":        detector.NewFeatureBranchDetector(),
		"NewCrissCrossMergeDetect":        detector.NewBranchMatrixDetector(detector.NewCrissCrossMergeDetect()),

		// Disabled
		// "NewFeatureBranchNameDetect": detector.NewBranchCompareDetector(detector.NewFeatureBranchNameDetect()),
		// "TwoParentsCommitDetect":     detector.NewCommitDetector(detector.TwoParentsCommitDetect()),
	}
	cacheDetectorRegistry = map[string]detector.CacheDetector{
		"ForcePushDetect": detector.NewCommitCacheDetector(detector.ForcePushDetect()),
	}
)

type Workflow struct {
	Name                    string
	WeightedCommitDetectors []WeightedDetector      `json:"-"`
	WeightedCacheDetectors  []WeightedCacheDetector `json:"-"`
	Violations              []violation.Violation   `json:"-"`
	Count                   int
	Total                   int
}

type WeightedDetector struct {
	Weight   int
	Detector detector.Detector
}

type WeightedCacheDetector struct {
	Weight   int
	Detector detector.CacheDetector
}

func GithubFlowWorkflow(cfg *config.Config) *Workflow {
	weightedCommitDetectors, weightedCacheDetectors := configureDetectors(cfg)

	return &Workflow{
		Name:                    "Github Flow",
		WeightedCommitDetectors: weightedCommitDetectors,
		WeightedCacheDetectors:  weightedCacheDetectors,
	}
}

// TODO: Remove this & use the config file instead. But it's currently useful for testing probably.
// LocalDetectors are detectors that can run locally without GitHub API calls.
func LocalDetectors() []detector.Detector {
	return []detector.Detector{
		detector.NewBranchDetector(detector.StaleBranchDetect()),
		detector.NewCommitDetector(detector.DiffMatchesMessageDetect()),
		detector.NewCommitDetector(detector.ShortCommitMessageDetect()),
		detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation()),
		detector.NewBranchCompareDetector(detector.NewBranchNameConsistencyDetect()),
		detector.NewCommitDetector(detector.BranchCommitDetect()), // used to check if branches are used
		detector.NewFeatureBranchDetector(),
		detector.NewBranchMatrixDetector(detector.NewCrissCrossMergeDetect()),
	}
}

// Run analysis on the git project for all the detectors defined by the workflow.
func (w *Workflow) Analyze(
	model *enriched.EnrichedModel,
	authors utils.Authors,
	current *cache.Cache,
	caches []*cache.Cache,
) (violated int,
	count,
	total int,
	violations []violation.Violation,
	err error,
) {
	v, c, t, vs, err := w.RunWeightedDetectors(model)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
	}
	add(&violated, &count, &total, &violations, v, c, t, &vs, 1)

	// Only run when we have a cache
	if current != nil || caches != nil {
		// assumes irst commit is the current user
		email := model.Commits[0].Committer.Email
		v, c, t, vs, err = w.RunCacheDetectors(model.Owner, model.Name, email, current, caches)
		if err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
		}
		add(&violated, &count, &total, &violations, v, c, t, &vs, 1)
	} else {
		log.Println("No cache loaded, skipping cache detectors")
	}

	w.Violations = violations
	w.Count = count
	w.Total = total

	return violated, count, total, violations, nil
}

func (w *Workflow) Result() (
	count,
	total int,
	violations []violation.Violation,
) {
	return w.Count, w.Total, w.Violations
}

func (w *Workflow) RunWeightedDetectors(model *enriched.EnrichedModel) (
	violated,
	count,
	total int,
	violations []violation.Violation,
	err error,
) {
	for _, wd := range w.WeightedCommitDetectors {
		if err := wd.Detector.Run(model); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to run weighted detectors: %w", err)
		}

		v, c, t, vs := wd.Detector.Result()
		add(&violated, &count, &total, &violations, v, c, t, &vs, wd.Weight)
	}

	return
}

// All cache detectors share the same current and cache, treated as readonly.
func (w *Workflow) RunCacheDetectors(owner, repo, email string, current *cache.Cache, caches []*cache.Cache) (
	int,
	int,
	int,
	[]violation.Violation,
	error,
) {
	violated, count, total := 0, 0, 0
	violations := []violation.Violation{}
	for _, wd := range w.WeightedCacheDetectors {
		if err := wd.Detector.Run(owner, repo, email, current, caches); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze caches: %w", err)
		}

		v, c, t, vs := wd.Detector.Result()
		add(&violated, &count, &total, &violations, v, c, t, &vs, wd.Weight)
	}

	// No violations means we can reset cache to current, otherwise append to cache
	var nc []*cache.Cache
	if len(violations) == 0 {
		nc = []*cache.Cache{current}
	} else {
		nc = append(nc, caches...)
		nc = append(nc, current)
	}
	if err := cache.WriteCaches(nc); err != nil {
		return 0, 0, 0, nil, fmt.Errorf("Failed to write cache: %w", err)
	}

	return violated, count, total, violations, nil
}

// Add weighted result a detector to the shared result.
func add(
	violated, count, total *int, violations *[]violation.Violation,
	v, c, t int, vs *[]violation.Violation,
	weight int,
) {
	*violated += v * weight
	*count += c * weight
	*total += t * weight
	*violations = append(*violations, *vs...)
}

// Summarize the results of the analysis into a csv file.
func (wk *Workflow) Csv(path, name, url string) error {
	exists, err := utils.Exists(path)
	if err != nil {
		return fmt.Errorf("Could not check if file exists: %w'", err)
	}

	var fh *os.File
	// nolint: gosec
	fh, err = os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("Failed to create csv file: %w", err)
	}

	w := csv.NewWriter(fh)
	// nolint: errcheck, gosec
	defer fh.Close()

	// Create file and header
	if !exists {
		header := []string{"Repository", "URL"}

		// Add detector names to header.
		for _, wd := range wk.WeightedCommitDetectors {
			t := reflect.TypeOf(wd.Detector)
			header = append(header, t.Elem().Name())
		}

		for _, wcd := range wk.WeightedCacheDetectors {
			t := reflect.TypeOf(wcd.Detector)
			header = append(header, t.Elem().Name())
		}

		err = w.Write(header)
		w.Flush()
		if err != nil {
			return fmt.Errorf("Could not write header to csv file: %w", err)
		}
	}

	// XXX: Body will likely change later with beyond length of headers, so don't alloc.
	// nolint: prealloc
	var body []string
	body = append(body, name, url)

	for _, wd := range wk.WeightedCommitDetectors {
		violated, _, _, _ := wd.Detector.Result()
		body = append(body, strconv.Itoa(violated))
	}

	for _, wcd := range wk.WeightedCacheDetectors {
		violated, _, _, _ := wcd.Detector.Result()
		body = append(body, strconv.Itoa(violated))
	}

	err = w.Write(body)
	if err != nil {
		return fmt.Errorf("could not write body to csv file: %w", err)
	}

	w.Flush()

	return nil
}

// Enable or disable detectors based on config.
func configureDetectors(cfg *config.Config) ([]WeightedDetector, []WeightedCacheDetector) {
	var weightedCommitDetectors []WeightedDetector
	var weightedCacheDetectors []WeightedCacheDetector

	for k := range cfg.Detectors {
		// Check keys match between config and registry.
		found := false
		if val, ok := detectorRegistry[k]; ok {
			weightedCommitDetectors = append(weightedCommitDetectors, WeightedDetector{
				Detector: val,
				Weight:   cfg.Detectors[k].Weight,
			})
			found = true
		}

		if val, ok := cacheDetectorRegistry[k]; ok {
			if cfg.Detectors[k].Enabled {
				weightedCacheDetectors = append(weightedCacheDetectors, WeightedCacheDetector{
					Detector: val,
					Weight:   cfg.Detectors[k].Weight,
				})
				found = true
			}
		}

		if !found {
			log.Printf("Detector \"%s\" from config not found", k)
		}
	}

	return weightedCommitDetectors, weightedCacheDetectors
}

// Write a JSON log of the workflow run.
// Assumes that the workflow is in a state that it has run to create a meaningful log.
func (w *Workflow) WriteLog(em enriched.EnrichedModel, cfg *config.Config) error {
	// Interface types within workflwo mean we need to reconsume the interface to get the concrete type.
	type LogViolation struct {
		Name         string
		Message      string
		Suggestion   string
		Email        string
		Author       remote.Author
		FileLocation string
		LineLocation int
		Severity     int
	}

	type Log struct {
		Date       time.Time
		Workflow   Workflow
		Config     config.Config
		Violations []LogViolation
		Model      enriched.EnrichedModel
	}

	LogViolations := make([]LogViolation, len(w.Violations))

	for i, v := range w.Violations {
		suggestion, _ := v.Suggestion()
		author, err := v.Author()
		if err != nil {
			author = &remote.Author{}
		}

		fileLocation, _ := v.FileLocation()
		lineLocation, _ := v.LineLocation()
		LogViolations[i] = LogViolation{
			Name:         v.Name(),
			Message:      v.Message(),
			Suggestion:   suggestion,
			Email:        v.Email(),
			Author:       *author,
			FileLocation: fileLocation,
			LineLocation: lineLocation,
			Severity:     int(v.Severity()),
		}
	}

	log := Log{
		Date:       time.Now(),
		Workflow:   *w,
		Config:     *cfg,
		Violations: LogViolations,
		Model:      em,
	}

	bytes, err := json.MarshalIndent(log, "", "")
	if err != nil {
		return fmt.Errorf("Failed to marshal workflow log: %w", err)
	}

	path := fmt.Sprintf("log-%s-%d.json", em.Name, time.Now().Unix())
	if err := ioutil.WriteFile(filepath.Clean(path), bytes, 0o600); err != nil {
		return fmt.Errorf("Error writing log to file: %w", err)
	}

	return nil
}
