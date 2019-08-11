package config

// Host mackerel check host info
type Host struct {
	// Manual definition
	ID            string            `json:"id" default:"hostname"`
	Hostname      string            `json:"hostname"`
	SourceType    string            `json:"sourceType" default:"host"`
	TargetRegion  string            `json:"targetRegion"`
	AssumeRoleARN string            `json:"assumeRoleArn"`
	Memos         map[string]string `json:"memos"`        // '{"ruleName":"memo",...}'
	CheckTargets  map[string]string `json:"checkTargets"` // '{"ruleName":"targetArn",...}'

	// Automatic definition
	HostID   string `json:"hostId"`
	CehckSum string `json:"checkSum"`
}

// CheckRule rule of check Plugin
type CheckRule struct {
	// Manual definition
	RuleName              string `json:"ruleName"`
	PluginType            string `json:"pluginType"`
	Parameters            string `json:"paramerters"`
	TimeoutSec            string `json:"timeoutSeconds"`
	PreventAlertAutoClose bool   `json:"preventAlertAutoClose"`
	CheckInterval         int    `json:"checkInterval"`
	Action                string `json:"action"`
	NotificationInterval  int    `json:"notificationInterval"` // <- Post report
	MaxCheckAttempts      int    `json:"maxCheckAttempts"`     // <- Post report

}
