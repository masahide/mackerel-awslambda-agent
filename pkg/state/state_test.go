package state

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamodb struct {
	dynamodbiface.DynamoDBAPI
	outputs []map[string]*dynamodb.AttributeValue
}

func (m *mockDynamodb) ScanPages(in *dynamodb.ScanInput, fn func(*dynamodb.ScanOutput, bool) bool) error {
	fn(&dynamodb.ScanOutput{Items: m.outputs}, true)
	return nil
}

func (m *mockDynamodb) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

type mockHostStore struct {
	key string
	res HostState
}

func (m *mockHostStore) ScanAll(out interface{}) error { return nil }
func (m *mockHostStore) Get(key string, out interface{}) error {
	m.key = key
	return nil
}

func (m *mockHostStore) Put(in interface{}) error {
	var ok bool
	if m.res, ok = in.(HostState); !ok {
		return errors.New("unmatched HostState type")
	}
	return nil
}

func TestPutHostState(t *testing.T) {
	tests := []struct {
		org      string
		hostname string
		expect   string
	}{
		{"orgname", "hostname", "orgname-hostname"},
		{"", "", "-"},
	}
	for n, tt := range tests {
		ms := &mockHostStore{}
		m := &Manager{
			TTLDays:  90,
			Org:      tt.org,
			Hostname: tt.hostname,
			Store:    ms,
		}
		err := m.PutHostState(HostState{})
		if err != nil {
			t.Error(err)
		}
		if ms.res.ID != tt.expect {
			t.Errorf("%d ID:%s expect:%s", n, ms.res.ID, tt.expect)
		}
	}
}

func TestGetHostState(t *testing.T) {
	tests := []struct {
		org      string
		hostname string
		expect   string
	}{
		{"orgname", "hostname", "orgname-hostname"},
	}
	for n, tt := range tests {
		ms := &mockHostStore{}
		m := &Manager{
			TTLDays:  90,
			Org:      tt.org,
			Hostname: tt.hostname,
			Store:    ms,
		}
		_, err := m.GetHostState()
		if err != nil {
			t.Error(err)
		}
		if ms.key != tt.expect {
			t.Errorf("%d ID:%s expect:%s", n, ms.key, tt.expect)
		}
	}
}

type mockCheckStore struct {
	key string
	res CheckState
}

func (m *mockCheckStore) ScanAll(out interface{}) error { return nil }
func (m *mockCheckStore) Get(key string, out interface{}) error {
	m.key = key
	ps, err := encodeCheckState(&CheckState{ID: key}, 0)
	if err != nil {
		return err
	}
	switch v := out.(type) {
	case *PluginState:
		*v = *ps
	}
	return nil
}

func (m *mockCheckStore) Put(in interface{}) error {
	ps, ok := in.(*PluginState)
	if !ok {
		return errors.New("unmatched PluginState type")
	}
	cs, err := decodeCheckState(ps)
	if err != nil {
		return err
	}
	m.res = *cs
	return nil
}

func TestPutCheckState(t *testing.T) {
	tests := []struct {
		name      string
		org       string
		hostname  string
		state     CheckState
		expectPut string
		expectGet string
	}{
		{"hoge", "org", "hostname", CheckState{"id000000000", []byte(`"json":"sample"`), ""}, "id000000000", "org-hostname-hoge"},
	}
	for n, tt := range tests {
		ms := &mockCheckStore{}
		m := &Manager{
			TTLDays:  90,
			Org:      tt.org,
			Hostname: tt.hostname,
			Store:    ms,
		}
		err := m.PutCheckState(tt.name, &tt.state)
		if err != nil {
			t.Error(err)
		}
		if ms.res.ID != tt.expectPut {
			t.Errorf("%d ID:%s expectPut:%s", n, ms.res.ID, tt.expectPut)
		}
		res, err := m.GetCheckState(tt.name)
		if err != nil {
			t.Error(err)
		}
		if res.ID != tt.expectGet {
			t.Errorf("%d ID:%s expectGet:%s", n, res.ID, tt.expectGet)
		}
	}
}
