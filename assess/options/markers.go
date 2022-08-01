package options

type Markers struct {
	Enable     []string `mapstructure:"enable"`
	Disable    []string `mapstructure:"disable"`
	EnableAll  bool     `mapstructure:"enable-all"`
	DisableAll bool     `mapstructure:"disable-all"`
	Fast       bool

	Presets []string
}
