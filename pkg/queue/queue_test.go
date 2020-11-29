package queue

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/masahide/mackerel-awslambda-agent/pkg/mock"
)

//go:generate mockgen  -destination ../mock/sqs.go -package mock github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI

func TestPostCheckReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := mock.NewMockSQSAPI(ctrl)
	mockSQS.EXPECT().
		SendMessageWithContext(gomock.Any(), gomock.Any()).
		Return(&sqs.SendMessageOutput{}, nil).
		Times(1)

	q, err := New(session.Must(session.NewSession()))
	if err != nil {
		t.Error(err)
	}
	q.svc = mockSQS
	err = q.PostCheckReport(context.Background(), &mackerel.CheckReport{})
	if err != nil {
		t.Error(err)
	}
}
