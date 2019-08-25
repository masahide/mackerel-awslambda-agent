package config

// CheckPluginParams is plugin awslambda event params
type CheckPluginParams struct {
	Org  string
	Name string
	CheckRule
	Host
}
