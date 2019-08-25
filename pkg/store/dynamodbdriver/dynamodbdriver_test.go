package dynamodbdriver

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func TestNew(t *testing.T) {
	test := []struct {
		tableName string
	}{
		{tableName: "hoge"},
	}
	for _, tt := range test {
		sess := session.Must(session.NewSession())
		d := New(sess, tt.tableName)
		if d.TableName != tt.tableName {
			t.Errorf("d.TableName=<%v> want <%v>", d.TableName, tt.tableName)
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
func (m *mockDynamodb) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.outputs = append(m.outputs, in.Item)
	return &dynamodb.PutItemOutput{}, nil
}

func (m *mockDynamodb) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":       {S: aws.String("test1")},
			"hostname": {S: aws.String("hostname1")},
			"checks": {L: []*dynamodb.AttributeValue{
				{M: map[string]*dynamodb.AttributeValue{"Name": {S: aws.String("check1")}, "Memo": {S: aws.String("fuga")}}},
			},
			},
		},
	}, nil
}

type host struct {
	ID       string
	Hostname string
	Checks   []check
}
type check struct {
	Name string
	Memo string
}

func TestScanAll(t *testing.T) {
	test := []struct {
		outputs []map[string]*dynamodb.AttributeValue
		want    []host
	}{
		{
			outputs: []map[string]*dynamodb.AttributeValue{
				{
					"id":       {S: aws.String("test1")},
					"hostname": {S: aws.String("hostname1")},
					"checks": {L: []*dynamodb.AttributeValue{
						{M: map[string]*dynamodb.AttributeValue{"Name": {S: aws.String("check1")}, "Memo": {S: aws.String("fuga")}}},
					}},
				},
				{
					"id":       {S: aws.String("test2")},
					"hostname": {S: aws.String("hostname2")},
					"checks": {L: []*dynamodb.AttributeValue{
						{M: map[string]*dynamodb.AttributeValue{"Name": {S: aws.String("check1")}, "Memo": {S: aws.String("fuga")}}},
					}},
				},
			},
			want: []host{
				{
					ID:       "test1",
					Hostname: "hostname1",
					Checks: []check{
						{Name: "check1", Memo: "fuga"},
					},
				},
				{
					ID:       "test2",
					Hostname: "hostname2",
					Checks: []check{
						{Name: "check1", Memo: "fuga"},
					},
				},
			},
		},
	}
	for _, tt := range test {
		m := mockDynamodb{
			outputs: tt.outputs,
		}
		d := &DynamoDB{DynamoDBAPI: &m}
		h := []host{}
		err := d.ScanAll(h)
		if err != nil {
			t.Error(err)
		}
		for i := range h {
			if !reflect.DeepEqual(h[i], tt.want[i]) {
				t.Errorf("h[i]=<%q> want <%q>", h[i], tt.want[i])
			}
		}
	}
}

func TestGet(t *testing.T) {
	test := []struct {
		key  string
		want string
	}{
		{
			key:  "test1",
			want: "test1",
		},
	}
	for _, tt := range test {
		m := mockDynamodb{}
		d := &DynamoDB{DynamoDBAPI: &m}
		h := host{}
		err := d.Get(tt.key, &h)
		if err != nil {
			t.Error(err)
		}
		if h.ID != tt.want {
			t.Errorf("h.ID=<%q> want <%q>", h.ID, tt.want)
		}
	}
}
func TestPut(t *testing.T) {
	m := mockDynamodb{}
	d := &DynamoDB{DynamoDBAPI: &m}
	h := host{
		ID:       "test1",
		Hostname: "hostname1",
		Checks: []check{
			{Name: "check1", Memo: "fuga"},
		},
	}
	err := d.Put(&h)
	if err != nil {
		t.Error(err)
	}
	want := []map[string]*dynamodb.AttributeValue{
		{
			"ID":       {S: aws.String("test1")},
			"Hostname": {S: aws.String("hostname1")},
			"Checks": {L: []*dynamodb.AttributeValue{
				{M: map[string]*dynamodb.AttributeValue{"Name": {S: aws.String("check1")}, "Memo": {S: aws.String("fuga")}}},
			}},
		},
	}
	if !reflect.DeepEqual(m.outputs[0], want[0]) {
		t.Errorf("m.Item=<%v> want <%v>", m.outputs, want)
	}
}
