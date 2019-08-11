package config

import (
	"os"
	"reflect"
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
			envs: []string{"HostsTable", "StateTable"},
			want: Env{
				HostsTable:      "mackerel-awslambda-hosts",
				CheckRulesTable: "mackerel-awslambda-checkrules",
				StateTable:      "mackerel-awslambda-state",
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
			t.Errorf("result = <%q> want <%q>", c.Env, tt.want)
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

func TestGetHosts(t *testing.T) {
	test := []struct {
		outputs []map[string]*dynamodb.AttributeValue
		want    map[string]Host
	}{
		{
			outputs: []map[string]*dynamodb.AttributeValue{
				{
					"id":       {S: aws.String("test1")},
					"hostname": {S: aws.String("hostname1")},
					"memos": {M: map[string]*dynamodb.AttributeValue{
						"check1": {S: aws.String("arn::hoge")},
					}},
				},
				{
					"id":       {S: aws.String("test2")},
					"hostname": {S: aws.String("hostname2")},
					"memos": {M: map[string]*dynamodb.AttributeValue{
						"check1": {S: aws.String("arn::hoge")},
					}},
				},
			},
			want: map[string]Host{
				"test1": {
					ID:       "test1",
					Hostname: "hostname1",
					Memos: map[string]string{
						"check1": "arn::hoge",
					},
				},
				"test2": {
					ID:       "test2",
					Hostname: "hostname2",
					Memos: map[string]string{
						"check1": "arn::hoge",
					},
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
		res, err := c.getHosts()
		if err != nil {
			t.Error(err)
		}
		for i := range res {
			if !reflect.DeepEqual(res[i], tt.want[i]) {
				t.Errorf("res[i]=<%q> want <%q>", res[i], tt.want[i])
			}
		}
	}
}

func TestGetCheckRules(t *testing.T) {
	test := []struct {
		outputs []map[string]*dynamodb.AttributeValue
		want    map[string]CheckRule
	}{
		{
			outputs: []map[string]*dynamodb.AttributeValue{
				{
					"ruleName":   {S: aws.String("test1")},
					"pluginType": {S: aws.String("hostname1")},
				},
				{
					"ruleName":   {S: aws.String("test2")},
					"pluginType": {S: aws.String("hostname2")},
				},
			},
			want: map[string]CheckRule{
				"test1": {
					RuleName:   "test1",
					PluginType: "cloudwatchlogs",
				},
				"test2": {
					RuleName:   "test2",
					PluginType: "cloudwatchlogs",
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
		res, err := c.getCheckRules()
		if err != nil {
			t.Error(err)
		}
		for i := range res {
			if !reflect.DeepEqual(res[i], tt.want[i]) {
				t.Errorf("res[i]=<%#v> want <%#v>", res[i], tt.want[i])
			}
		}
	}
}

func TestGetCheckStates(t *testing.T) {
	test := []struct {
		outputs []map[string]*dynamodb.AttributeValue
		want    map[string]CheckState
	}{
		{
			outputs: []map[string]*dynamodb.AttributeValue{
				{
					"stateID": {S: aws.String("test1")},
					"data":    {B: []byte("hostname1")},
				},
				{
					"stateID": {S: aws.String("test2")},
					"data":    {B: []byte("hostname1")},
				},
			},
			want: map[string]CheckState{
				"test1": {
					StateID: "test1",
					Data:    []byte("hostname1"),
				},
				"test2": {
					StateID: "test2",
					Data:    []byte("hostname1"),
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
		res, err := c.getCheckStates()
		if err != nil {
			t.Error(err)
		}
		for i := range res {
			if !reflect.DeepEqual(res[i], tt.want[i]) {
				t.Errorf("res[i]=<%q> want <%q>", res[i], tt.want[i])
			}
		}
	}
}
