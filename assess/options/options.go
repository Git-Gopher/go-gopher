package options

type Options struct {
	optionsDir string // The directory containing the options.
	Run        Run

	CutoffDate       string `mapstructure:"cutoff-date"` // viper doesn't support time.Time, parse manual for now
	DefaultAlgorithm string `mapstructure:"default-algorithm"`
	OutputDir        string `mapstructure:"output-dir"`
	OutputRepoFolder bool   `mapstructure:"output-repo-folder"`
	LoginWhiteList   string `mapstructure:"login-whitelist"`

	FilenameTemplate string `mapstructure:"filename-template"`
	HeaderTemplate   string `mapstructure:"header-template"`

	MarkersSettings MarkersSettings `mapstructure:"markers-settings"`
	Markers         Markers         `mapstructure:"markers"`
}

func (o *Options) GetOptionsDir() string {
	return o.optionsDir
}

func NewDefault() *Options {
	return &Options{
		LoginWhiteList:  "github-classroom[bot]",
		MarkersSettings: defaultMarkersSettings,
	}
}

type Version struct {
	Format string `mapstructure:"format"`
}
