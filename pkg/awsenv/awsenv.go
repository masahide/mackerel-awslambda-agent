package awsenv

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"golang.org/x/xerrors"
)

type Credential struct {
	AccessKey       string `json:"aws_access_key_id" toml:"aws_access_key_id"`
	SecretAccessKey string `json:"aws_secret_access_key" toml:"aws_secret_access_key"`
	SessionToken    string `json:"aws_session_token" toml:"aws_session_token"`
}

type Credentials map[string]Credential

func getAWSEnvs() Credential {
	return Credential{
		AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	}
}

func unsetCredential() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
}

func putCredsFile(home string, creds Credentials) error {
	data, err := toml.Marshal(creds)
	if err != nil {
		return xerrors.Errorf("toml.Marshal err:%w", err)
	}

	return ioutil.WriteFile(filepath.Join(home, ".aws", "credentials"), data, 0600)
}

/*
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
*/

func readCreds(home string) Credentials {
	data, err := ioutil.ReadFile(filepath.Join(home, ".aws", "credentials"))
	if err != nil {
		log.Printf("readFile credentials err:%s", err)

		return map[string]Credential{}
	}
	res := Credentials{}
	err = toml.Unmarshal(data, &res)
	if err != nil {
		log.Printf("credentials Unmarshal err:%s", err)

		return map[string]Credential{}
	}

	return res
}

func EnvToCredentialFile(profile, home string) error {
	/*
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
	*/
	cred := getAWSEnvs()
	if len(cred.AccessKey) == 0 {
		return nil
	}
	// nolint:errcheck
	os.Mkdir(filepath.Join(home, ".aws"), 0755)
	creds := readCreds(home)
	creds[profile] = cred
	if err := putCredsFile(home, creds); err != nil {
		return err
	}
	unsetCredential()

	return nil
}
