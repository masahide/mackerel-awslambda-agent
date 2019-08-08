package config

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func TestNewAgentConfig(t *testing.T) {

	test := []struct {
		envs []string
		want Env
	}{
		{
			envs: []string{"CheckConfigTable", "StateTable"},
			want: Env{
				CheckConfigTable: "mackerel-serverless-check-config",
				StateTable:       "mackerel-serverless-state",
			},
		},
	}
	sess := session.Must(session.NewSession())
	for _, tt := range test {
		for _, k := range tt.envs {
			os.Unsetenv(k)
		}
		c := NewAgentConfig(sess)
		if c.Env != tt.want {
			t.Errorf("result = <%s> want <%s>", c, tt.want)
		}
	}
}

type mockDynamodb struct {
	dynamodbiface.DynamoDBAPI
	outputs []map[string]*dynamodb.AttributeValue
}

func (m *mockDynamodb) ScanPages(in *dynamodb.ScanInput, fn func(*dynamodb.ScanOutput, bool) bool) error {
	fn(&dynamodb.ScanOutput{Items: m.outputs}, true)
	return nil
}

func TestGetCheckConfigs(t *testing.T) {
	test := []struct {
		outputs []map[string]*dynamodb.AttributeValue
		want    []CheckConfig
	}{
		{
			outputs: []map[string]*dynamodb.AttributeValue{
				map[string]*dynamodb.AttributeValue{
					"Id":       &dynamodb.AttributeValue{S: aws.String("test1")},
					"Hostname": &dynamodb.AttributeValue{S: aws.String("hostname1")},
					"name":     &dynamodb.AttributeValue{S: aws.String("name1")},
				},
				map[string]*dynamodb.AttributeValue{
					"Id":       &dynamodb.AttributeValue{S: aws.String("test2")},
					"Hostname": &dynamodb.AttributeValue{S: aws.String("hostname2")},
					"name":     &dynamodb.AttributeValue{S: aws.String("name2")},
				},
			},
			want: []CheckConfig{
				CheckConfig{
					ID:       "test1",
					Hostname: "hostname1",
					Name:     "name1",
				},
				CheckConfig{
					ID:       "test2",
					Hostname: "hostname2",
					Name:     "name2",
				},
			},
		},
	}
	sess := session.Must(session.NewSession())
	for _, tt := range test {
		c := NewAgentConfig(sess)
		m := mockDynamodb{
			outputs: tt.outputs,
		}
		c.DynamoDBAPI = &m
		res, err := c.GetCheckConfigs()
		if err != nil {
			t.Error(err)
		}
		for i := range res {
			if res[i] != tt.want[i] {
				t.Errorf("res[i]=<%q> want <%q>", res[i], tt.want[i])
			}
		}
	}
}
