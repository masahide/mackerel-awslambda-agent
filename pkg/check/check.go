package check

import (
	"context"
	"log"
	"strings"
	"time"

	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mackerel-container-agent/cmdutil"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
)

const (
	stateFileKeyword = "%{STATE}"
)

// Check check plugin struct
type Check struct {
	config.Host
	config.CheckRule
	state.CheckState
	hostID string
}

func (c *Check) replaceStateFilePath(path string) string {
	return strings.ReplaceAll(c.Command, stateFileKeyword, path)
}

// Generate generates check report
func (c *Check) Generate(ctx context.Context) (*mackerel.CheckReport, error) {
	cmd := cmdutil.CommandString(c.Command)
	now := time.Now()
	stdout, stderr, exitCode, err := cmdutil.RunCommand(ctx, cmd, "", c.Env, c.Timeout)

	if stderr != "" {
		log.Printf("plugin %s (%s): %q", c.Name, c.Command, stderr)
	}

	var message string
	var status mackerel.CheckStatus
	if err != nil {
		log.Printf("Warning plugin %s (%s): %s", c.Name, c.Command, err)
		message = err.Error()
		status = mackerel.CheckStatusUnknown
	} else {
		message = stdout
		status = exitCodeToStatus(exitCode)
	}

	newReport := mackerel.CheckReport{
		Source:               mackerel.NewCheckSourceHost(c.hostID),
		Name:                 c.Name,
		Status:               status,
		Message:              message,
		OccurredAt:           now.Unix(),
		NotificationInterval: c.NotificationInterval,
		MaxCheckAttempts:     c.MaxCheckAttempts,
	}

	LatestReport := c.LatestReport
	c.LatestReport = &newReport
	if LatestReport == nil {
		return &newReport, nil
	}
	if LatestReport.Status == mackerel.CheckStatusOK && newReport.Status == mackerel.CheckStatusOK {
		// do not report ok -> ok
		return nil, nil
	}
	return &newReport, nil
}

var exitCodeMap = map[int]mackerel.CheckStatus{
	0: mackerel.CheckStatusOK,
	1: mackerel.CheckStatusWarning,
	2: mackerel.CheckStatusCritical,
}

func exitCodeToStatus(exitCode int) mackerel.CheckStatus {
	if code, ok := exitCodeMap[exitCode]; ok {
		return code
	}
	return mackerel.CheckStatusUnknown
}
