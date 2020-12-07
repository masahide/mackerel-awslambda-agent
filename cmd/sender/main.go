package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	lambdaevent "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/sender"
)

var (
	s *sender.Sender
)

func main() {
	var env config.Env
	var sess *session.Session
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	sess = session.Must(session.NewSession())
	conf, err := config.LoadS3Config(sess, env.S3Bucket, env.S3Key)
	if err != nil {
		log.Fatal(err.Error())
	}
	s = sender.New(conf.Apikey)
	lambdaevent.Start(handler)
}

// Handler aws lambda handler.
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	return s.Run(sqsEvent)
}
