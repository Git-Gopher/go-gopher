package options

var defaultMarkersSettings = MarkersSettings{}

type MarkersSettings struct {
	Commit        CommitSettings       `mapstructure:"commit"`
	CommitMessage CommitMessageSetting `mapstructure:"commit-message"`
	Branching     BranchingSettings    `mapstructure:"branching"`
	PullRequest   PullRequestSettings  `mapstructure:"pull-request"`
	General       GeneralSettings      `mapstructure:"general"`
}

type ThresholdSettings struct {
	ThresholdA int `mapstructure:"threshold-a"`
	ThresholdB int `mapstructure:"threshold-b"`
	ThresholdC int `mapstructure:"threshold-c"`
}

type CommitSettings struct {
	GradingAlgorithm  string             `mapstructure:"grading-algorithm"`
	ThresholdSettings *ThresholdSettings `mapstructure:"threshold-settings"`
}

type CommitMessageSetting struct {
	GradingAlgorithm  string             `mapstructure:"grading-algorithm"`
	ThresholdSettings *ThresholdSettings `mapstructure:"threshold-settings"`
}

type BranchingSettings struct {
	GradingAlgorithm  string             `mapstructure:"grading-algorithm"`
	ThresholdSettings *ThresholdSettings `mapstructure:"threshold-settings"`
}

type PullRequestSettings struct {
	GradingAlgorithm  string             `mapstructure:"grading-algorithm"`
	ThresholdSettings *ThresholdSettings `mapstructure:"threshold-settings"`
}

type GeneralSettings struct {
	GradingAlgorithm  string             `mapstructure:"grading-algorithm"`
	ThresholdSettings *ThresholdSettings `mapstructure:"threshold-settings"`
}
