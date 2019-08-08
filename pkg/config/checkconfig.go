package config

// CheckConfig mackerel check plugin config
type CheckConfig struct {
	ID                   string `json:"Id" default:"Name-hostname-pluginType"`
	Name                 string `json:"name"`
	SourceType           string `json:"sourceType" default:"host"`
	Hostname             string `json:"hostname"`
	HostID               string `json:"hostId"`
	Message              string `json:"message"`
	NotificationInterval int    `json:"notificationInterval"`
	MaxCheckAttempts     int    `json:"maxCheckAttempts"`
	PluginType           string `json:"pluginType"`
	Parameters           string `json:"paramerters"`
}
