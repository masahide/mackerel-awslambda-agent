package awsenv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvToCredentialFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up

	os.Setenv("AWS_ACCESS_KEY_ID", "keyxxxxxxxxxxxxxxx")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secretxxxxxxxxxxxxxxx")
	os.Setenv("AWS_SESSION_TOKEN", "tokenxxxxxxxxxxxxxxx")
	if err = EnvToCredentialFile("test_profile", dir); err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadFile(filepath.Join(dir, ".aws", "credentials"))
	if err != nil {
		t.Error(err)
	}
	expected := `[test_profile]
aws_access_key_id     = keyxxxxxxxxxxxxxxx
aws_secret_access_key = secretxxxxxxxxxxxxxxx
aws_session_token     = tokenxxxxxxxxxxxxxxx
`
	if diff := cmp.Diff(string(data), expected); diff != "" {
		t.Errorf("credentials file differs: (-got +want)\n%s", diff)
	}

	os.Setenv("AWS_ACCESS_KEY_ID", "keyxxxxxxxxxxxxxxx")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secretxxxxxxxxxxxxxxx")
	os.Setenv("AWS_SESSION_TOKEN", "tokenxxxxxxxxxxxxxxx")
	if err = EnvToCredentialFile("test2_profile", dir); err != nil {
		t.Error(err)
	}
	data, err = ioutil.ReadFile(filepath.Join(dir, ".aws", "credentials"))
	if err != nil {
		t.Error(err)
	}
	expected = `[test_profile]
aws_access_key_id     = keyxxxxxxxxxxxxxxx
aws_secret_access_key = secretxxxxxxxxxxxxxxx
aws_session_token     = tokenxxxxxxxxxxxxxxx

[test2_profile]
aws_access_key_id     = keyxxxxxxxxxxxxxxx
aws_secret_access_key = secretxxxxxxxxxxxxxxx
aws_session_token     = tokenxxxxxxxxxxxxxxx
`
	if diff := cmp.Diff(string(data), expected); diff != "" {
		t.Errorf("credentials file differs: (-got +want)\n%s", diff)
	}
}
