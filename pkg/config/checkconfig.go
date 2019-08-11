package config

// Host mackerel check host info
type Host struct {
	ID            string            `json:"id" default:"hostname"` // Primary key
	Hostname      string            `json:"hostname"`
	SourceType    string            `json:"sourceType" default:"host"`
	TargetRegion  string            `json:"targetRegion"`
	AssumeRoleARN string            `json:"assumeRoleArn"`
	Memos         map[string]string `json:"memos"`        // '{"ruleName":"memo",...}'
	CheckTargets  map[string]string `json:"checkTargets"` // '{"ruleName":"targetArn",...}'

}

// CheckRule rule of check Plugin
type CheckRule struct {
	RuleName              string `json:"ruleName"` // Primary key
	PluginType            string `json:"pluginType"`
	Parameters            string `json:"paramerters"`
	TimeoutSec            int    `json:"timeoutSeconds"`
	PreventAlertAutoClose bool   `json:"preventAlertAutoClose"`
	CheckInterval         int    `json:"checkInterval"`
	Action                string `json:"action"`
	NotificationInterval  int    `json:"notificationInterval"` // <- Post report
	MaxCheckAttempts      int    `json:"maxCheckAttempts"`     // <- Post report

}

// State of check Plugin & host
type State struct {
	// Manual definition
	ID           string `json:"id"`    // Primary key
	State        []byte `json:"state"` // State data
	HostID       string `json:"hostId"`
	HostCehckSum string `json:"HostCheckSum"`
}
