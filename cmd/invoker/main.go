package main

import (
	"context"
	"log"

	lambdaevent "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/invoker"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store/dynamodbdriver"
)

var (
	agent *config.AgentConfig
	sess  *session.Session
	env   config.Env
	iv    *invoker.Invoker
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	sess = session.Must(session.NewSession())
	s := dynamodbdriver.New(sess, env.StateTable)
	agent, err = config.NewAgentConfig(s, sess)
	if err != nil {
		log.Fatal(err)
	}
	iv = invoker.New(sess, env)
	lambdaevent.Start(handler)
}

// Handler aws lambda handler.
func handler(ctx context.Context, event config.Env) error {
	if err := agent.LoadAgentConfig(sess, env.S3Bucket, env.S3Key); err != nil {
		log.Printf("agent.LoadAgentConfig err: %+v", err)
		// nolint:wrapcheck
		return err
	}
	if err := agent.GetHost(); err != nil {
		log.Printf("agent.GetHost err: %+v", err)
		// nolint:wrapcheck
		return err
	}
	if err := iv.Run(agent); err != nil {
		log.Printf("invoker Run err: %+v", err)
		// nolint:wrapcheck
		return err
	}
	return nil
}
