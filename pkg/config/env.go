package config

// Env config of mackerel-awslambda-agent
type Env struct {
	HostsTable      string `default:"mackerel-awslambda-hosts"`
	CheckRulesTable string `default:"mackerel-awslambda-checkrules"`
	StateTable      string `default:"mackerel-awslambda-state"`
	StateTTLDays    int    `default:"90"`
	TempDir         string `default:""`
}
