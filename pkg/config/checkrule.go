package config

import "time"

// CheckRule rule of check Plugin.
type CheckRule struct {
	Name                  string        `json:"name"` // Primary key
	PluginType            string        `json:"pluginType"`
	Command               string        `json:"command"`
	Env                   []string      `json:"env"`
	Timeout               time.Duration `json:"timeout"`
	PreventAlertAutoClose bool          `json:"preventAlertAutoClose"`
	CheckInterval         int           `json:"checkInterval"`
	Action                string        `json:"action"`
	NotificationInterval  uint          `json:"notificationInterval"` // <- Post report
	MaxCheckAttempts      uint          `json:"maxCheckAttempts"`     // <- Post report
}
