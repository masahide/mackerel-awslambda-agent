package config

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/kelseyhightower/envconfig"
)

// Env config of mackerel-awslambda-agent
type Env struct {
	HostsTable      string `default:"mackerel-awslambda-hosts"`
	CheckRulesTable string `default:"mackerel-awslambda-checkrules"`
	StateTable      string `default:"mackerel-awslambda-state"`
}

// AgentConfig is agent config struct
type AgentConfig struct {
	Env
	dynamodbiface.DynamoDBAPI
	CheckRules map[string]CheckRule
	Hosts      map[string]Host
}

// NewAgentConfig load config from env
func NewAgentConfig(p client.ConfigProvider) *AgentConfig {
	a := &AgentConfig{
		DynamoDBAPI: dynamodb.New(p),
		CheckRules:  map[string]CheckRule{},
		Hosts:       map[string]Host{},
	}
	err := envconfig.Process("", &a.Env)
	if err != nil {
		log.Fatal(err.Error())
	}
	return a
}

// getHosts get check configs
func (a *AgentConfig) getHosts() (map[string]Host, error) {
	hosts := []Host{}
	var unmarshalErr error
	err := a.ScanPages(
		&dynamodb.ScanInput{TableName: &a.HostsTable},
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			h := make([]Host, len(page.Items))
			if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &h); unmarshalErr != nil {
				return false
			}
			hosts = append(hosts, h...)
			return true
		},
	)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	if err != nil {
		return nil, err
	}
	res := make(map[string]Host, len(hosts))
	for _, h := range hosts {
		res[h.ID] = h
	}
	return res, nil
}

func (a *AgentConfig) getCheckRules() (map[string]CheckRule, error) {
	rules := []CheckRule{}
	var unmarshalErr error
	err := a.ScanPages(
		&dynamodb.ScanInput{TableName: &a.CheckRulesTable},
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			c := make([]CheckRule, len(page.Items))
			if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &c); unmarshalErr != nil {
				return false
			}
			rules = append(rules, c...)
			return true
		},
	)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	if err != nil {
		return nil, err
	}
	res := make(map[string]CheckRule, len(rules))
	for _, c := range rules {
		res[c.RuleName] = c
	}
	return res, nil
}

func (a *AgentConfig) getCheckStates() (map[string]CheckState, error) {
	states := []CheckState{}
	var unmarshalErr error
	err := a.ScanPages(
		&dynamodb.ScanInput{TableName: &a.StateTable},
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			c := make([]CheckState, len(page.Items))
			if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &c); unmarshalErr != nil {
				return false
			}
			states = append(states, c...)
			return true
		},
	)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	if err != nil {
		return nil, err
	}
	res := make(map[string]CheckState, len(states))
	for _, c := range states {
		res[c.StateID] = c
	}
	return res, nil
}
