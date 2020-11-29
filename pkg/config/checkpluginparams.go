package config

import "github.com/masahide/mackerel-awslambda-agent/pkg/state"

// CheckPluginParams is plugin awslambda event params.
type CheckPluginParams struct {
	Org   string          `json:"org"`
	Name  string          `json:"name"`
	Rule  CheckRule       `json:"rule"`
	Host  Host            `json:"host"`
	State state.HostState `json:"state"`
}
