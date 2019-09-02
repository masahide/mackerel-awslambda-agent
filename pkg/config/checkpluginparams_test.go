package config

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
)

func TestJSONMarshal(t *testing.T) {
	params := CheckPluginParams{
		Org:  "org",
		Name: "name",
		Rule: CheckRule{
			Name:                  "ruleName",
			PluginType:            "type",
			Command:               "command option1 op2",
			Env:                   []string{"aaa", "bbbb"},
			Timeout:               1 * time.Second,
			PreventAlertAutoClose: false,
			CheckInterval:         1,
			Action:                "action",
			NotificationInterval:  1,
			MaxCheckAttempts:      1,
		},
		Host:  Host{},
		State: state.HostState{},
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}
	want := []byte("aaa")
	if bytes.Compare(b, want) != 0 {
		t.Errorf("%s, want:%s", b, want)
	}
}
