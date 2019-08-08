package config

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/kelseyhightower/envconfig"
)

// Env config of mackerel-serverless-agent
type Env struct {
	CheckConfigTable string `default:"mackerel-serverless-check-config"`
	StateTable       string `default:"mackerel-serverless-state"`
}

// AgentConfig is agent config struct
type AgentConfig struct {
	Env
	dynamodbiface.DynamoDBAPI
}

// NewAgentConfig load config from env
func NewAgentConfig(p client.ConfigProvider) *AgentConfig {
	a := &AgentConfig{
		DynamoDBAPI: dynamodb.New(p),
	}
	err := envconfig.Process("", &a.Env)
	if err != nil {
		log.Fatal(err.Error())
	}
	return a
}

// GetCheckConfigs get check configs
func (a *AgentConfig) GetCheckConfigs() ([]CheckConfig, error) {
	res := []CheckConfig{}
	var unmarshalErr error
	err := a.ScanPages(&dynamodb.ScanInput{
		TableName: &a.CheckConfigTable,
	}, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		conf := make([]CheckConfig, len(page.Items))
		if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &conf); unmarshalErr != nil {
			return false
		}
		res = append(res, conf...)
		return true
	})
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return res, err
}
