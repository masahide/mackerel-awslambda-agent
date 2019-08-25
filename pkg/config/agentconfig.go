package config

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/kelseyhightower/envconfig"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store/dynamodbdriver"
)

// AgentConfig is agent config struct
type AgentConfig struct {
	Env
	CheckRules     map[string]CheckRule
	Hosts          map[string]Host
	hostsStore     store.Store
	checkRuleStore store.Store
	stateStore     store.Store
}

// NewAgentConfig load config from env
func NewAgentConfig(p client.ConfigProvider) *AgentConfig {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	a := &AgentConfig{
		hostsStore:     dynamodbdriver.New(p, env.HostsTable),
		checkRuleStore: dynamodbdriver.New(p, env.CheckRulesTable),
		stateStore:     dynamodbdriver.New(p, env.StateTable),
		CheckRules:     map[string]CheckRule{},
		Hosts:          map[string]Host{},
		Env:            env,
	}
	return a
}

// LoadTables load config from env
func (a *AgentConfig) LoadTables() error {
	var err error
	if a.Hosts, err = a.getHosts(); err != nil {
		return err
	}
	if a.CheckRules, err = a.getCheckRules(); err != nil {
		return err
	}
	//a.CheckStates, err = a.getStates()
	return err
}

// getHosts get check configs
func (a *AgentConfig) getHosts() (map[string]Host, error) {
	var hosts []Host
	if err := a.hostsStore.ScanAll(&hosts); err != nil {
		return nil, err
	}
	res := make(map[string]Host, len(hosts))
	for _, h := range hosts {
		res[h.ID] = h
	}
	return res, nil
}

func (a *AgentConfig) getCheckRules() (map[string]CheckRule, error) {
	var rules []CheckRule
	if err := a.checkRuleStore.ScanAll(&rules); err != nil {
		return nil, err
	}
	res := make(map[string]CheckRule, len(rules))
	for _, c := range rules {
		res[c.Name] = c
	}
	return res, nil
}

/*
func (a *AgentConfig) getStates() (map[string]state.CheckState, error) {
		var states []state.CheckState

		err := a.ScanPages(
			&dynamodb.ScanInput{TableName: &a.StateTable},
			func(page *dynamodb.ScanOutput, lastPage bool) bool {
				c := make([]state.CheckState, len(page.Items))
				if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &c); unmarshalErr != nil {
					return falsecheckRuleStore
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
		res := make(map[string]state.CheckState, len(states))
		for _, c := range states {
			res[c.ID] = c
		}
		return res, nil
	return nil, nil
}

// PutState put state.CheckState
func (a *AgentConfig) PutState(in state.CheckState) error {
	attr, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return errors.Wrap(err, "MarshalMap err")
	}
	_, err = a.PutItem(&dynamodb.PutItemInput{Item: attr})
	return err
}
*/
