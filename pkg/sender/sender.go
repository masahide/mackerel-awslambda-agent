package sender

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mackerelio/mackerel-client-go"
	"golang.org/x/xerrors"
)

type Sender struct {
	apiKey string
}

type reportData struct {
	Source struct {
		Type   string `json:"type"`
		HostID string `json:"hostId"`
	} `json:"source"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	OccurredAt int64  `json:"occurredAt"`
}

func New(apiKey string) *Sender {
	return &Sender{
		apiKey: apiKey,
	}
}
func (s *Sender) Run(sqsEvent events.SQSEvent) error {
	checkReports := make([]*mackerel.CheckReport, 0, len(sqsEvent.Records))
	for _, record := range sqsEvent.Records {
		report := reportData{}
		err := json.Unmarshal([]byte(record.Body), &report)
		if err != nil {
			return xerrors.Errorf("json.Unmarshal data:%s err: %w", record.Body, err)
		}

		checkReports = append(checkReports, &mackerel.CheckReport{
			Source:     mackerel.NewCheckSourceHost(report.Source.HostID),
			Name:       report.Name,
			Status:     mackerel.CheckStatus(report.Status),
			Message:    report.Message,
			OccurredAt: report.OccurredAt,
		})
	}
	client := mackerel.NewClient(s.apiKey)
	if err := client.PostCheckReports(&mackerel.CheckReports{Reports: checkReports}); err != nil {
		return xerrors.Errorf("mackerel PostCheckReports data:%# v err: %w", checkReports, err)
	}
	return nil
}
