package config

// Host mackerel check host info
type Host struct {
	ID            string  `json:"id" default:"hostname"` // Primary key
	Hostname      string  `json:"hostname"`
	SourceType    string  `json:"sourceType" default:"host"`
	TargetRegion  string  `json:"targetRegion"`
	AssumeRoleARN string  `json:"assumeRoleArn"`
	Checks        []Check `json:"checks" dynamodbav:"checks"` // '[{"name":"checkName","memo":""},...]'
}

// Check mackerel check info
type Check struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// CheckRule rule of check Plugin
type CheckRule struct {
	Name                  string `json:"name"` // Primary key
	PluginType            string `json:"pluginType"`
	Command               string `json:"command"`
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
	HostCehckSum string `json:"hostCheckSum"`
}

/*
func (u Checks) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &h); unmarshalErr != nil { if av.M == nil {
		return nil
	}

	n, err := strconv.ParseInt(*av.N, 10, 0)
	if err != nil {
		return err
	}

	u.Value = int(n)
	return nil
}
*/
