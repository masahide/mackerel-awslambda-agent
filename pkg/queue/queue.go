package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/kelseyhightower/envconfig"
	"github.com/mackerelio/mackerel-client-go"
	"golang.org/x/xerrors"
)

type sqsParams struct {
	CheckReportQueueARN string
}

type Queue struct {
	envs sqsParams
	svc  sqsiface.SQSAPI
}

func New(sess *session.Session) (*Queue, error) {
	q := Queue{
		svc: sqs.New(sess),
	}
	if err := envconfig.Process("", &q.envs); err != nil {
		return nil, err
	}
	return &q, nil
}

func (q *Queue) PostCheckReport(ctx context.Context, report *mackerel.CheckReport) error {
	data, err := json.Marshal(report)
	if err != nil {
		return xerrors.Errorf("json.Marshal err:%s", err)
	}
	_, err = q.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		QueueUrl:    &q.envs.CheckReportQueueARN,
		MessageBody: aws.String(string(data)),
	})
	return err
}
