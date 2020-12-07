package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/kelseyhightower/envconfig"
	"github.com/mackerelio/mackerel-client-go"
	"golang.org/x/xerrors"
)

type sqsParams struct {
	QueueURL string
}

type Queue struct {
	envs sqsParams
	svc  sqsiface.SQSAPI
}

// New Queue struct.
func New(sess client.ConfigProvider) (*Queue, error) {
	q := Queue{
		svc: sqs.New(sess),
	}
	if err := envconfig.Process("", &q.envs); err != nil {
		return nil, xerrors.Errorf("envconfig.Process err: %w", err)
	}

	return &q, nil
}

// PostCheckReport sendMessage to sqs queue.
func (q *Queue) PostCheckReport(ctx context.Context, report *mackerel.CheckReport) error {
	data, err := json.Marshal(report)
	if err != nil {
		return xerrors.Errorf("json.Marshal err: %w", err)
	}
	// nolint:exhaustivestruct
	_, err = q.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		QueueUrl:    &q.envs.QueueURL,
		MessageBody: aws.String(string(data)),
	})
	if err != nil {
		return xerrors.Errorf("SendMessageWithContext ARN:%s err: %w", q.envs.QueueURL, err)
	}

	return nil
}
