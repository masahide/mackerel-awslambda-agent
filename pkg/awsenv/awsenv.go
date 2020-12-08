package awsenv

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"golang.org/x/xerrors"
)

const (
	awsAccessKeyID = "aws_access_key_id"
	// nolint:gosec
	awsSecretAccessKey = "aws_secret_access_key"
	awsSessionToken    = "aws_session_token"

	envAwsAccessKeyID = "AWS_ACCESS_KEY_ID"
	// nolint:gosec
	envAwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	// nolint:gosec
	envAwsSessionToken = "AWS_SESSION_TOKEN"

	envAwsAccessKey = "AWS_ACCESS_KEY"
	// nolint:gosec
	envAwsSecretKey     = "AWS_SECRET_KEY"
	envAwsSecurityToken = "AWS_SECURITY_TOKEN"
)

type Credential struct {
	AccessKey       string
	SecretAccessKey string
	SessionToken    string
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getAWSEnvs() Credential {
	return Credential{
		AccessKey:       os.Getenv(envAwsAccessKeyID),
		SecretAccessKey: os.Getenv(envAwsSecretAccessKey),
		SessionToken:    os.Getenv(envAwsSessionToken),
	}
}

func unsetCredential() {
	os.Unsetenv(envAwsAccessKeyID)
	os.Unsetenv(envAwsAccessKey)
	os.Unsetenv(envAwsSecretAccessKey)
	os.Unsetenv(envAwsSecretKey)
	os.Unsetenv(envAwsSecurityToken)
	os.Unsetenv(envAwsSessionToken)
}

func putCredsFile(credPath, profile string, cred Credential) error {
	config, err := ini.Load(credPath)
	if err != nil {
		return xerrors.Errorf("ini.Load err: %w", err)
	}
	section := config.Section(profile)
	section.Key(awsAccessKeyID).SetValue(cred.AccessKey)
	section.Key(awsSecretAccessKey).SetValue(cred.SecretAccessKey)
	section.Key(awsSessionToken).SetValue(cred.SessionToken)
	err = config.SaveTo(credPath)
	if err != nil {
		return xerrors.Errorf("config.SaveTo err: %w", err)
	}
	return nil
}

func EnvToCredentialFile(profile, home string) error {
	credPath := filepath.Join(home, ".aws", "credentials")
	if !exists(credPath) {
		// nolint:errcheck
		os.Mkdir(filepath.Join(home, ".aws"), 0755)
		file, err := os.Create(credPath)
		if err != nil {
			log.Fatal(err)
			return xerrors.Errorf("create file err: %w", err)
		}
		file.Close()
	}
	cred := getAWSEnvs()
	if len(cred.AccessKey) == 0 {
		return nil
	}
	if err := putCredsFile(credPath, profile, cred); err != nil {
		return err
	}
	unsetCredential()

	return nil
}
