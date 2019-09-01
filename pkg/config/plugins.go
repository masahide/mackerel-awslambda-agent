package config

import "github.com/masahide/mackerel-awslambda-agent/pkg/state"

// CheckPluginParams is plugin awslambda event params
type CheckPluginParams struct {
	Org  string
	Name string
	CheckRule
	Host
	State state.HostState
}
