package invoker

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"golang.org/x/xerrors"
)

type Invoker struct {
	lambdaSvc lambdaiface.LambdaAPI
	env       config.Env
}

func New(sess *session.Session, env config.Env) *Invoker {
	svc := lambda.New(sess)
	return &Invoker{
		lambdaSvc: svc,
		env:       env,
	}
}
func (iv *Invoker) Run(agent *config.AgentConfig) error {
	for _, checker := range agent.CheckRules {
		checkConf := config.CheckPluginParams{
			//Org:       agent.HostState.Organization,
			Rule:      checker,
			HostState: agent.HostState,
		}
		payload, err := json.Marshal(checkConf)
		if err != nil {
			return xerrors.Errorf("CheckerPluginParams json.Marshal err: %w", err)
		}
		_, err = iv.lambdaSvc.Invoke(&lambda.InvokeInput{
			FunctionName:   &iv.env.CheckerFunc,
			Payload:        payload,
			InvocationType: aws.String(lambda.InvocationTypeEvent),
		})
		if err != nil {
			return xerrors.Errorf("lambda invoke err: %w", err)
		}
		log.Printf("invoke lambda param: %s", payload)
	}
	return nil
}
