package statefile

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestPutStateFiles(t *testing.T) {
	test := []struct {
		input []byte
		want  []byte
	}{
		{
			input: []byte(`{"11111/aaa":"Y2NjY2NjCg==","11111/bb":"Y2NjY2NjZWZlZgo=","11111/bbbb":"Y2NjY2NjZWZlZmFhCg==","fuga":"YWZmYWRmYQo="}`),
			want:  []byte(`{"11111/aaa":"Y2NjY2NjCg==","11111/bb":"Y2NjY2NjZWZlZgo=","11111/bbbb":"Y2NjY2NjZWZlZmFhCg==","fuga":"YWZmYWRmYQo="}`),
		},
		{
			input: []byte(`{}`),
			want:  []byte(`{}`),
		},
	}
	for i, tt := range test {
		dir, err := ioutil.TempDir("", "mackerel-awslambda-agent-statefile")
		if err != nil {
			t.Fatal(err)
		}
		if err := PutStatefiles(dir, tt.input); err != nil {
			t.Error(err)
		}
		b, err := GetStatefiles(dir)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(b, tt.want) {
			t.Errorf("i=%d b=%s,want=%s", i, b, tt.want)
		}
		os.RemoveAll(dir) // clean up
	}
}
