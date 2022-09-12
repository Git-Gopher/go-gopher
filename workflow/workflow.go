package workflow

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/config"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	"github.com/Git-Gopher/go-gopher/violation"
	log "github.com/sirupsen/logrus"
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
		"BranchNameConsistencyDetect":     detector.NewBranchCompareDetector(detector.BranchNameConsistencyDetect()),
		"FeatureBranchDetector":           detector.NewFeatureBranchDetector("FeatureBranchDetector"),
		"CrissCrossMergeDetect":           detector.NewBranchMatrixDetector(detector.CrissCrossMergeDetect()),
		"UnresolvedDetect":                detector.NewCommitDetector(detector.UnresolvedDetect()),
		"EmptyCommitDetect":               detector.NewCommitDetector(detector.EmptyCommitDetect()),
		"BinaryDetect":                    detector.NewCommitDetector(detector.BinaryDetect()),

		// Disabled
		// "NewFeatureBranchNameDetect": detector.NewBranchCompareDetector(detector.NewFeatureBranchNameDetect()),
		// "TwoParentsCommitDetect":     detector.NewCommitDetector(detector.TwoParentsCommitDetect()),
	}
	cacheDetectorRegistry = map[string]detector.CacheDetector{
		"ForcePushDetect": detector.NewCommitCacheDetector(detector.ForcePushDetect()),
	}
)

type Workflow struct {
	Name                    string                  `json:"name"`
	WeightedCommitDetectors []WeightedDetector      `json:"-"`
	WeightedCacheDetectors  []WeightedCacheDetector `json:"-"`
	Violations              []violation.Violation   `json:"-"`
	Count                   int                     `json:"count"`
	Total                   int                     `json:"total"`
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
		detector.NewCommitDetector(detector.UnresolvedDetect()),
		detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation()),
		detector.NewBranchCompareDetector(detector.BranchNameConsistencyDetect()),
		detector.NewCommitDetector(detector.BranchCommitDetect()), // used to check if branches are used
		detector.NewFeatureBranchDetector("FeatureBranchDetector"),
		detector.NewBranchMatrixDetector(detector.CrissCrossMergeDetect()),
	}
}

// Run analysis on the git project for all the detectors defined by the workflow.
func (w *Workflow) Analyze(
	model *enriched.EnrichedModel,
	authors utils.Authors,
	current *cache.Cache,
	previous []*cache.Cache,
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
	if current != nil || previous != nil {
		// assumes irst commit is the current user
		email := model.Commits[0].Committer.Email
		v, c, t, vs, err = w.RunCacheDetectors(model.Owner, model.Name, email, current, previous)
		if err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
		}
		add(&violated, &count, &total, &violations, v, c, t, &vs, 1)
	} else {
		log.Println("No cache loaded, skipping cache detectors...")
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
func (w *Workflow) RunCacheDetectors(owner, repo, email string, current *cache.Cache, previous []*cache.Cache) (
	int,
	int,
	int,
	[]violation.Violation,
	error,
) {
	violated, count, total := 0, 0, 0
	violations := []violation.Violation{}
	for _, wd := range w.WeightedCacheDetectors {
		if err := wd.Detector.Run(owner, repo, email, current, previous); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("failed to analyze caches: %w", err)
		}

		v, c, t, vs := wd.Detector.Result()
		add(&violated, &count, &total, &violations, v, c, t, &vs, wd.Weight)
	}

	var nc []*cache.Cache
	nc = append(nc, previous...)
	nc = append(nc, current)

	if err := cache.Write(nc); err != nil {
		return 0, 0, 0, nil, fmt.Errorf("failed to write cache: %w", err)
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
	//nolint: gosec
	fh, err = os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("Failed to create csv file: %w", err)
	}

	w := csv.NewWriter(fh)
	//nolint: errcheck, gosec
	defer fh.Close()

	// Create file and header
	if !exists {
		header := []string{"Repository", "URL"}

		// Add detector names to header.
		for _, wd := range wk.WeightedCommitDetectors {
			header = append(header, wd.Detector.Name())
		}

		for _, wcd := range wk.WeightedCacheDetectors {
			header = append(header, wcd.Detector.Name())
		}

		err = w.Write(header)
		w.Flush()
		if err != nil {
			return fmt.Errorf("Could not write header to csv file: %w", err)
		}
	}

	// XXX: Body will likely change later with beyond length of headers, so don't alloc.
	//nolint: prealloc
	var body []string
	body = append(body, name, url)

	for _, wd := range wk.WeightedCommitDetectors {
		violated, found, total, _ := wd.Detector.Result()

		body = append(body, fmt.Sprintf("%d:%d:%d", violated, found, total))
	}

	for _, wcd := range wk.WeightedCacheDetectors {
		violated, found, total, _ := wcd.Detector.Result()
		body = append(body, fmt.Sprintf("%d:%d:%d", violated, found, total))
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
func (w *Workflow) WriteLog(em enriched.EnrichedModel, cfg *config.Config) (string, error) {
	// Interface types within workflow mean we need to reconsume the interface to get the concrete type.
	type logViolation struct {
		Name         string
		Message      string
		Suggestion   string
		Email        string
		Time         string
		Author       remote.Author
		FileLocation string
		LineLocation int
		Severity     int
	}

	type log struct {
		Name       string         `json:"name"`
		URL        string         `json:"url"`
		Date       time.Time      `json:"date"`
		Violations []logViolation `json:"violations"`
		Workflow   Workflow       `json:"workflow"`
		Config     config.Config  `json:"config"`
	}

	LogViolations := make([]logViolation, len(w.Violations))

	for i, v := range w.Violations {
		suggestion, _ := v.Suggestion()
		author, err := v.Author()
		if err != nil {
			author = &remote.Author{}
		}

		fileLocation, _ := v.FileLocation()
		lineLocation, _ := v.LineLocation()
		LogViolations[i] = logViolation{
			Name:         v.Name(),
			Message:      v.Message(),
			Suggestion:   suggestion,
			Email:        v.Email(),
			Time:         v.Time().Format(time.UnixDate),
			Author:       *author,
			FileLocation: fileLocation,
			LineLocation: lineLocation,
			Severity:     int(v.Severity()),
		}
	}

	l := log{
		Name:       em.Name,
		URL:        em.URL,
		Date:       time.Now(),
		Violations: LogViolations,
		Workflow:   *w,
		Config:     *cfg,
	}

	bytes, err := json.MarshalIndent(l, "", "")
	if err != nil {
		return "", fmt.Errorf("Failed to marshal workflow log: %w", err)
	}

	fn := fmt.Sprintf("/home/wqsz7xn/Projects/go-gopher/output2/log-%s-%d.json", em.Name, time.Now().Unix())
	if err := os.WriteFile(filepath.Clean(fn), bytes, 0o600); err != nil {
		return "", fmt.Errorf("failed writing log to file: %w", err)
	}

	return fn, nil
}

// Print a summary of the workflow violations to stdout.
func PrintSummary(authors utils.Authors, v, c, t int, vs []violation.Violation) {
	var violations, suggestions []violation.Violation
	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			violations = append(violations, v)
		case violation.Suggestion:
			suggestions = append(suggestions, v)
		}
	}

	var vSd strings.Builder
	for _, v := range violations {
		vSd.WriteString(v.Display(authors))
	}
	markup.Group("Violations", vSd.String())

	var sSd strings.Builder
	for _, v := range suggestions {
		sSd.WriteString(v.Display(authors))
	}
	markup.Group("Suggestions", sSd.String())

	var aSd strings.Builder
	counts := make(map[string]int)
	for _, v := range vs {
		email := v.Email()
		login, err := authors.Find(email)
		if err != nil {
			continue
		}
		counts[*login]++
	}

	for login, count := range counts {
		aSd.WriteString(fmt.Sprintf("%s: %d\n", login, count))
	}

	aSd.WriteString(fmt.Sprintf("violated: %d\n", v))
	aSd.WriteString(fmt.Sprintf("count: %d\n", c))
	aSd.WriteString(fmt.Sprintf("total: %d\n", t))
	markup.Group("Summary", aSd.String())
}

// Create a markdown summary for a workflow, inluding a summary of the violations and suggestions.
// Usually used in pull request comments.
func MarkdownSummary(authors utils.Authors, vs []violation.Violation) string { // nolint: gocognit
	md := markup.CreateMarkdown("Workflow Summary")
	md.AddLine(fmt.Sprintf("Created with git-gopher version `%s`", version.BuildVersion()))

	// Separate violation types.
	var violations []violation.Violation
	var suggestions []violation.Violation

	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			if v.Current() {
				violations = append(violations, v)
			}
		case violation.Suggestion:
			if v.Current() {
				suggestions = append(suggestions, v)
			}
		default:
			log.Printf("Unknown violation severity: %v", v.Severity())
		}
	}

	if len(violations) > 0 {
		headers := []string{"Violation", "Message", "Advice", "Author"}
		rows := make([][]string, len(violations))

		for i, v := range violations {
			row := make([]string, len(headers))
			name := v.Name()
			row[0] = name
			message := v.Message()
			row[1] = message

			suggestion, err := v.Suggestion()
			if err != nil {
				suggestion = ""
			}
			row[2] = suggestion

			usernamePtr, err := authors.Find(v.Email())
			if err != nil || usernamePtr == nil {
				row[3] = "unknown"
			} else {
				row[3] = markup.Author(*usernamePtr).Markdown()
			}

			rows[i] = row
		}

		md.BeginCollapsable("Violations")
		md.Table(headers, rows)
		md.EndCollapsable()
	}

	if len(suggestions) > 0 {
		headers := []string{"Suggestion", "Message", "Advice", "Author"}
		rows := make([][]string, len(suggestions))

		for i, v := range suggestions {
			row := make([]string, len(headers))
			name := v.Name()
			row[0] = name
			message := v.Message()
			row[1] = message

			suggestion, err := v.Suggestion()
			if err != nil {
				suggestion = ""
			}
			row[2] = suggestion

			usernamePtr, err := authors.Find(v.Email())
			if err != nil || usernamePtr == nil {
				row[3] = "unknown"
			} else {
				row[3] = markup.Author(*usernamePtr).Markdown()
			}

			rows[i] = row
		}

		md.BeginCollapsable("Suggestions")
		md.Table(headers, rows)
		md.EndCollapsable()
	}

	workflowUrl := os.Getenv("WORKFLOW_URL")
	if (len(violations)+len(suggestions)) < len(vs) && workflowUrl != "" {
		md.AddLine(fmt.Sprintf(`There still exist some violations beyond the scope of this pull request, 
			please view the full log [here](%s)`, workflowUrl))
	}

	// Google form
	md.AddLine(fmt.Sprintf("Have any feedback? Feel free to submit it [here](%s)", utils.GoogleFormURL))

	return md.Render()
}
