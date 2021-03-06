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
		Rule: CheckRule{
			Name:                  "ruleName",
			PluginType:            "type",
			Command:               "command option1 op2",
			Env:                   []string{"aaa", "bbbb"},
			Timeout:               1 * time.Second,
			PreventAlertAutoClose: false,
			CheckInterval:         1,
			//Action:                "action",
			NotificationInterval: 1,
			MaxCheckAttempts:     1,
		},
		HostState: &state.HostState{},
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}
	//want := []byte(`{"org":"org","name":"name","rule":{"name":"ruleName","pluginType":"type","command":"command option1 op2","env":["aaa","bbbb"],"timeout":1000000000,"preventAlertAutoClose":false,"checkInterval":1,"action":"action","notificationInterval":1,"maxCheckAttempts":1},"host":{"id":"","hostname":"","sourceType":"","targetRegion":"","assumeRoleArn":"","checks":null},"state":{"id":"","hostId":"","hostCheckSum":""}}`)
	//want := []byte(`{"org":"org","name":"name","rule":{"name":"ruleName","pluginType":"type","command":"command option1 op2","env":["aaa","bbbb"],"timeout":1000000000,"preventAlertAutoClose":false,"checkInterval":1,"memo":"","customIdentifier":"","notificationInterval":1,"maxCheckAttempts":1},"host":{"id":"","hostname":"","sourceType":"","targetRegion":"","assumeRoleArn":"","checks":null},"state":{"id":"","hostId":"","hostCheckSum":""}}`)
	want := []byte(`{"rule":{"name":"ruleName","pluginType":"type","command":"command option1 op2","env":["aaa","bbbb"],"timeout":1000000000,"preventAlertAutoClose":false,"checkInterval":1,"memo":"","customIdentifier":"","notificationInterval":1,"maxCheckAttempts":1},"state":{"id":"","organization":"","hostname":"","hostId":"","hostCheckSum":""}}`)
	if bytes.Compare(b, want) != 0 {
		t.Errorf("%s, want:%s", b, want)
	}
}
