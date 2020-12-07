package config

import "github.com/masahide/mackerel-awslambda-agent/pkg/state"

// CheckPluginParams is plugin awslambda event params.

type CheckPluginParams struct {
	//Org  string    `json:"org"`
	Rule CheckRule `json:"rule"`
	//Host      Host             `json:"host"`
	HostState *state.HostState `json:"state"`
}
