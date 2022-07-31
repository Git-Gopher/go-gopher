package options

type Run struct {
	IsVerbose   bool `mapstructure:"verbose"`
	Silent      bool
	Concurrency int
}
