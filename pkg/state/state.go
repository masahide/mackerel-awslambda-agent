package state

import (
	"encoding/json"
	"time"

	"github.com/masahide/mackerel-awslambda-agent/pkg/store"
	"golang.org/x/xerrors"
)

const (
	defaultCheckStateTTLDays = 90
)

// CheckState is check plugin state.
type CheckState struct {
	ID           string `json:"id"`
	StateFiles   []byte `json:"stateFiles,omitempty"` // JSON data
	LatestStatus string `json:"latestReport,omitempty"`
}

// HostState is dynamodb table of hostID.
type HostState struct {
	ID           string `json:"id"` // Primary key ( m.Org + "-" + m.Hostname
	Organization string `json:"organization"`
	Hostname     string `json:"hostname"`
	HostID       string `json:"hostId" dynamodbav:",omitempty"`
	HostCheckSum string `json:"hostCheckSum" dynamodbav:",omitempty"`
}

// PluginState is dynamodb table of check Plugin state.
type PluginState struct {
	ID    string `json:"id"`                            // Primary key (HostID:hostname
	State []byte `json:"state" dynamodbav:",omitempty"` // State data
	TTL   int    `json:"ttl" dynamodbav:",omitempty"`
}

// Manager struct.
type Manager struct {
	TTLDays  int
	Org      string
	Hostname string
	store.Store
}

func (m *Manager) ttl() int64 {
	ttl := m.TTLDays
	if ttl == 0 {
		ttl = defaultCheckStateTTLDays
	}

	return time.Now().AddDate(0, 0, ttl).Unix()
}

// GetCheckState Extract checkState from pluginState.
func (m *Manager) GetCheckState(name string) (*CheckState, error) {
	ps, err := m.GetPluginState(name)
	if err != nil {
		return nil, xerrors.Errorf("GetPluginState err: %w", err)
	}

	return decodeCheckState(ps)
}

func decodeCheckState(in *PluginState) (*CheckState, error) {
	var res CheckState
	err := json.Unmarshal(in.State, &res)
	if err != nil {
		return nil, xerrors.Errorf("decodeCheckState json.Unmarshal data:[%s] err: %w", in.State, err)
	}
	if len(res.StateFiles) == 0 {
		res.StateFiles = []byte("{}")
	}

	return &res, nil
}

// PutCheckState Extract checkState from pluginState.
func (m *Manager) PutCheckState(name string, in *CheckState) error {
	ps, err := encodeCheckState(in, m.ttl())
	if err != nil {
		return err
	}

	return m.PutPluginState(name, ps)
}

func encodeCheckState(in *CheckState, ttl int64) (*PluginState, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return nil, xerrors.Errorf("json.Marshal err: %w", err)
	}

	return &PluginState{
		ID:    in.ID,
		State: b,
		TTL:   int(ttl),
	}, nil
}

// GetHostState is get hostID from dynamodb.
func (m *Manager) GetHostState() (*HostState, error) {
	id := m.Org + "-" + m.Hostname
	res := &HostState{}
	if err := m.Get(id, res); err != nil {
		return nil, xerrors.Errorf("m.Get err: %w", err)
	}

	return res, nil
}

// PutHostState put  host id.
func (m *Manager) PutHostState(in HostState) error {
	in.ID = m.Org + "-" + m.Hostname

	return m.Put(in)
}

// GetPluginState is get plugin state from dynamodb.
func (m *Manager) GetPluginState(name string) (*PluginState, error) {
	id := m.Org + "-" + m.Hostname + "-" + name
	res := &PluginState{}
	if err := m.Get(id, res); err != nil {
		return nil, xerrors.Errorf("m.Get err: %w", err)
	}
	if len(res.ID) == 0 {
		res.ID = m.Org + "-" + m.Hostname + "-" + name
	}
	if len(res.State) == 0 {
		res.State = []byte("{}")
	}

	return res, nil
}

// PutPluginState put plugin state.
func (m *Manager) PutPluginState(name string, in *PluginState) error {
	in.ID = m.Org + "-" + m.Hostname + "-" + name

	return m.Put(in)
}
