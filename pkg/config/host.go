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
