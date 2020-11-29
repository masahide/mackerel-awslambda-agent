package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/aws/aws-lambda-go/events"
	lambdaevent "github.com/aws/aws-lambda-go/lambda"
	"github.com/mackerelio/mackerel-container-agent/cmdutil"
	"github.com/mackerelio/mackerel-container-agent/config"
)

func main() {
	lambdaevent.Start(handler)

}

// Handler aws lambda handler
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	env := config.Env([]string{})
	dir := os.Getenv("LAMBDA_TASK_ROOT")
	cmd := cmdutil.CommandString(fmt.Sprintf("%s -h", path.Join(dir, "check-aws-cloudwatch-logs")))
	//cmd := cmdutil.CommandString(fmt.Sprintf("set -x;pwd;ls -la %s", dir))
	stdout, stderr, exitCode, err := cmdutil.RunCommand(ctx, cmd, "", env, 60*time.Second)

	log.Printf("cmd:%s\nstdout:%s\n,stderr:%s\n,exitCode:%d\n,err:%s", cmd, stdout, stderr, exitCode, err)
	return err
}
