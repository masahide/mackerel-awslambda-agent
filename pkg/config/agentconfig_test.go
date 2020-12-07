package config

import (
	"os"
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

func TestLoadTables(t *testing.T) {
	sess := session.Must(session.NewSession())
	c := NewAgentConfig(sess)
	c.hostsStore = &mockCheckStore{checkOutputs: []CheckRule{}}
	if len(c.CheckRules) != 0 {
		t.Error("len(c.CheckRules)!=0")
	}
}
