package config

// Env config of mackerel-awslambda-agent.
type Env struct {
	StateTable   string `default:"mackerel-awslambda-state"`
	StateTTLDays int    `default:"90"`
	S3Key        string `default:"mackerel.toml"`
	S3Bucket     string
	Hostname     string `default:"hostname"`
	Organization string `default:"organization"`
	TempDir      string `default:""`
	CheckerFunc  string
}
