package checkplugin

import (
	"context"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
)

//checkplugin_test.go:47: res[i]=<&mackerel.CheckReport{Source:(*mackerel.checkSourceHost)( 0xc00000f420), Name:"", Status:"OK", Message:"", OccurredAt: 0, NotificationInterval: 0x0, MaxCheckAttempts: 0x0}> want <mackerel.CheckReport{Source:mackerel.CheckSource(nil), Name:"test", Status:"OK", Message:"", OccurredAt: 0, NotificationInterval: 0x0, MaxCheckAttempts: 0x0}>

func TestGenerate(t *testing.T) {
	test := []struct {
		params config.CheckPluginParams
		ps     state.PluginState
		want   mackerel.CheckReport
	}{
		{
			params: config.CheckPluginParams{
				HostState: &state.HostState{},
			},
			ps: state.PluginState{
				ID:    "1111111",
				State: []byte{},
				TTL:   0,
			},
			want: mackerel.CheckReport{
				Name:    "",
				Status:  "OK",
				Message: "",
			},
		},
	}
	for _, tt := range test {
		// config.CheckPluginParams
		m := mockHostStore{
			ps: tt.ps,
		}
		c := NewCheckPlugin(&m, tt.params)
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
