package main

import (
	"context"
	"log"
	"os"

	lambdaevent "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/masahide/mackerel-awslambda-agent/pkg/awsenv"
	"github.com/masahide/mackerel-awslambda-agent/pkg/checkplugin"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/queue"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store/dynamodbdriver"
	"golang.org/x/xerrors"
)

const (
	homeDir = "/tmp/home"
)

var (
	// nolint:gochecknoglobals
	// nolint:gochecknoglobals
	q *queue.Queue
	s store.Store
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var env config.Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	sess := session.Must(session.NewSession())
	s = dynamodbdriver.New(sess, env.StateTable)
	if q, err = queue.New(sess); err != nil {
		log.Fatal(err)
	}
	os.Mkdir(homeDir, 0755)
	os.Setenv("HOME", homeDir)
	if err := awsenv.EnvToCredentialFile("default", homeDir); err != nil {
		log.Fatal(err)
	}
	lambdaevent.Start(handler)
}

func handler(ctx context.Context, params config.CheckPluginParams) error {
	check := checkplugin.NewCheckPlugin(s, params)
	report, err := check.Generate(ctx)
	if err != nil {
		log.Printf("check.Generate err:%+v", err)

		return xerrors.Errorf("check.Generate err: %w", err)
	}
	if report == nil {
		return nil
	}
	if err := q.PostCheckReport(ctx, report); err != nil {
		log.Printf("PostCheckReport err:%+v", err)

		return xerrors.Errorf("PostCheckReport err: %w", err)
	}

	return nil
}
