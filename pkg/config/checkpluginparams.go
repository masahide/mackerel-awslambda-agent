package config

import "github.com/masahide/mackerel-awslambda-agent/pkg/state"

// CheckPluginParams is plugin awslambda event params.

type CheckPluginParams struct {
	Rule      CheckRule        `json:"rule"`
	HostState *state.HostState `json:"state"`
}
