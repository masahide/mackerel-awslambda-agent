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
					"pluginType": {S: aws.String("cloudwatchlogs")},
				},
				{
					"ruleName":   {S: aws.String("test2")},
					"pluginType": {S: aws.String("cloudwatchlogs")},
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
				t.Errorf("key:%s res[i]=<%#v> want <%#v>", i, res[i], tt.want[i])
			}
		}
	}
}

func TestGetStates(t *testing.T) {
	test := []struct {
		outputs []map[string]*dynamodb.AttributeValue
		want    map[string]State
	}{
		{
			outputs: []map[string]*dynamodb.AttributeValue{
				{
					"id":    {S: aws.String("test1")},
					"state": {B: []byte("hostname1")},
				},
				{
					"id":    {S: aws.String("test2")},
					"state": {B: []byte("hostname1")},
				},
			},
			want: map[string]State{
				"test1": {
					ID:    "test1",
					State: []byte("hostname1"),
				},
				"test2": {
					ID:    "test2",
					State: []byte("hostname1"),
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
		res, err := c.getStates()
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

func (m *mockDynamodb) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func TestPutState(t *testing.T) {
	sess := session.Must(session.NewSession())
	c := NewAgentConfig(sess)
	m := mockDynamodb{}
	c.DynamoDBAPI = &m
	err := c.PutState(State{})
	if err != nil {
		t.Error(err)
	}
}

func TestLoadTables(t *testing.T) {
	sess := session.Must(session.NewSession())
	c := NewAgentConfig(sess)
	m := mockDynamodb{
		outputs: []map[string]*dynamodb.AttributeValue{},
	}
	c.DynamoDBAPI = &m
	err := c.LoadTables()
	if err != nil {
		t.Error(err)
	}
	if len(c.Hosts) != 0 {
		t.Error("len(c.Hosts)!=0")
	}
	if len(c.CheckRules) != 0 {
		t.Error("len(c.CheckRules)!=0")
	}
	if len(c.States) != 0 {
		t.Error("len(c.States)!=0")
	}
}
