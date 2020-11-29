package main

import (
	"context"
	"log"
	"os"

	lambdaevent "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/masahide/mackerel-awslambda-agent/pkg/awsenv"
	"github.com/masahide/mackerel-awslambda-agent/pkg/checkplugin"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/queue"
)

var (
	sess *session.Session
	q    *queue.Queue
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if err := awsenv.EnvToCredentialFile("default", home); err != nil {
		log.Fatal(err)
	}
	sess = session.Must(session.NewSession())
	if q, err = queue.New(sess); err != nil {
		log.Fatal(err)
	}
	lambdaevent.Start(handler)
}

func handler(ctx context.Context, params config.CheckPluginParams) error {

	check := checkplugin.NewCheckPlugin(sess, params)
	report, err := check.Generate(ctx)
	if err != nil {
		log.Printf("check.Generate err:%v", err)
		return err
	}
	if report == nil {
		return nil
	}
	if err := q.PostCheckReport(ctx, report); err != nil {
		log.Printf("PostCheckReport err:%v", err)
		return err
	}
	return nil
}
