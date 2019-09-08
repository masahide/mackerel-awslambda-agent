package checkplugin

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
)

func TestGenerate(t *testing.T) {
	test := []struct {
		params config.CheckPluginParams
		ps     state.PluginState
		want   mackerel.CheckReport
	}{
		{
			params: config.CheckPluginParams{
				Name: "test",
			},
			ps: state.PluginState{
				ID:    "1111111",
				State: []byte{},
				TTL:   0,
			},
			want: mackerel.CheckReport{
				Name:    "test",
				Status:  "OK",
				Message: "",
			},
		},
	}
	sess := session.Must(session.NewSession())
	for _, tt := range test {
		//config.CheckPluginParams
		c := NewCheckPlugin(sess, tt.params)
		m := mockHostStore{
			ps: tt.ps,
		}
		c.Manager.Store = &m
		ctx := context.Background()
		res, err := c.Generate(ctx)
		if err != nil {
			t.Errorf("Generate err:%+v", err)
		}
		res.OccurredAt = 0
		if res.Name != tt.want.Name {
			t.Errorf("res[i]=<%# v> want <%# v>", res, tt.want)
		}
		if res.Status != tt.want.Status {
			t.Errorf("res[i]=<%# v> want <%# v>", res, tt.want)
		}
		if res.Message != tt.want.Message {
			t.Errorf("res[i]=<%# v> want <%# v>", res, tt.want)
		}
	}
}

type mockHostStore struct {
	ps state.PluginState
}

func (m *mockHostStore) ScanAll(out interface{}) error { return nil }
func (m *mockHostStore) Get(key string, out interface{}) error {
	switch v := out.(type) {
	case *state.PluginState:
		*v = m.ps
	}
	return nil
}
func (m *mockHostStore) Put(in interface{}) error { return nil }
