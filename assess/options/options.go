package options

type Options struct {
	optionsDir string // The directory containing the options.
	Run        Run

	DefaultAlgorithm string `mapstructure:"default-algorithm"`
	OutputDir        string `mapstructure:"output-dir"`
	OutputRepoFolder bool   `mapstructure:"output-repo-folder"`

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
		MarkersSettings: defaultMarkersSettings,
	}
}

type Version struct {
	Format string `mapstructure:"format"`
}
