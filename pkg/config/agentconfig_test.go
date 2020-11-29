package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
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
				StateTTLDays:    90,
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

type mockHostStore struct {
	hostOutputs []Host
}

func (m *mockHostStore) ScanAll(out interface{}) error {
	switch v := out.(type) {
	case *[]Host:
		*v = m.hostOutputs
	}
	return nil
}
func (m *mockHostStore) Get(key string, out interface{}) error { return nil }
func (m *mockHostStore) Put(in interface{}) error              { return nil }

func TestGetHosts(t *testing.T) {
	test := []struct {
		outputs []Host
		want    map[string]Host
	}{
		{
			outputs: []Host{
				{
					ID:       "test1",
					Hostname: "hostname1",
					Checks: []Check{
						{Name: "check1", Memo: "fuga"},
					},
				},
				{
					ID:       "test2",
					Hostname: "hostname2",
					Checks: []Check{
						{Name: "check1", Memo: "fuga"},
					},
				},
			},
			want: map[string]Host{
				"test1": {
					ID:       "test1",
					Hostname: "hostname1",
					Checks: []Check{
						{Name: "check1", Memo: "fuga"},
					},
				},
				"test2": {
					ID:       "test2",
					Hostname: "hostname2",
					Checks: []Check{
						{Name: "check1", Memo: "fuga"},
					},
				},
			},
		},
	}
	sess := session.Must(session.NewSession())
	for _, tt := range test {
		c := NewAgentConfig(sess)
		m := mockHostStore{
			hostOutputs: tt.outputs,
		}
		c.hostsStore = &m
		res, err := c.getHosts()
		if err != nil {
			t.Error(err)
		}
		for i := range tt.want {
			if !reflect.DeepEqual(res[i], tt.want[i]) {
				t.Errorf("res[i]=<%q> want <%q>", res[i], tt.want[i])
			}
		}
	}
}

type mockCheckStore struct {
	checkOutputs []CheckRule
}

func (m *mockCheckStore) ScanAll(out interface{}) error {
	switch v := out.(type) {
	case *[]CheckRule:
		*v = m.checkOutputs
	}
	return nil
}
func (m *mockCheckStore) Get(key string, out interface{}) error { return nil }
func (m *mockCheckStore) Put(in interface{}) error              { return nil }

func TestGetCheckRules(t *testing.T) {
	test := []struct {
		outputs []CheckRule
		want    map[string]CheckRule
	}{
		{
			outputs: []CheckRule{
				{
					Name:       "test1",
					PluginType: "cloudwatchlogs",
				},
				{
					Name:       "test2",
					PluginType: "cloudwatchlogs",
				},
			},
			want: map[string]CheckRule{
				"test1": {
					Name:       "test1",
					PluginType: "cloudwatchlogs",
				},
				"test2": {
					Name:       "test2",
					PluginType: "cloudwatchlogs",
				},
			},
		},
	}
	sess := session.Must(session.NewSession())
	for _, tt := range test {
		c := NewAgentConfig(sess)
		m := mockCheckStore{
			checkOutputs: tt.outputs,
		}
		c.checkRuleStore = &m
		c.hostsStore = &m
		res, err := c.getCheckRules()
		if err != nil {
			t.Error(err)
		}
		for i := range tt.want {
			if !reflect.DeepEqual(res[i], tt.want[i]) {
				t.Errorf("key:%s res[i]=<%#v> want <%#v>", i, res[i], tt.want[i])
			}
		}
	}
}

func TestLoadTables(t *testing.T) {
	sess := session.Must(session.NewSession())
	c := NewAgentConfig(sess)
	c.checkRuleStore = &mockHostStore{hostOutputs: []Host{}}
	c.hostsStore = &mockCheckStore{checkOutputs: []CheckRule{}}
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
}
