package checkplugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/kelseyhightower/envconfig"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mackerel-container-agent/cmdutil"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
	"github.com/masahide/mackerel-awslambda-agent/pkg/statefile"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store/dynamodbdriver"
	"golang.org/x/xerrors"
)

const (
	tempDirPrefix = "mackerel"
)

var exitCodeMap = map[int]mackerel.CheckStatus{
	0: mackerel.CheckStatusOK,
	1: mackerel.CheckStatusWarning,
	2: mackerel.CheckStatusCritical,
}

// CheckPlugin struct
type CheckPlugin struct {
	config.CheckPluginParams
	config.Env
	*state.Manager
	*state.CheckState
	tmpDir string
}

// NewCheckPlugin Create Plugin struct
func NewCheckPlugin(p client.ConfigProvider, params config.CheckPluginParams) *CheckPlugin {
	var env config.Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	plugin := &CheckPlugin{
		CheckPluginParams: params,
		Env:               env,
		Manager: &state.Manager{
			TTLDays:  env.StateTTLDays,
			Org:      params.Org,
			Hostname: params.Host.Hostname,
			Store:    dynamodbdriver.New(p, env.StateTable),
		},
	}
	return plugin
}

// Generate generates check report
func (c *CheckPlugin) Generate(ctx context.Context) (*mackerel.CheckReport, error) {
	if err := c.initialize(); err != nil {
		return nil, xerrors.Errorf("initialize err:%w", err)
	}
	defer func() {
		if err := c.finalize(); err != nil {
			log.Printf("finalize err:%s", err)
		}
	}()
	return c.generate(ctx)
}

// initialize is load config of CheckPlugin
func (c *CheckPlugin) initialize() error {
	var err error
	c.CheckState, err = c.GetCheckState(c.Rule.Name)
	if err != nil {
		return xerrors.Errorf("GetCheckState err:%w", err)
	}
	c.tmpDir, err = ioutil.TempDir(c.Env.TempDir, tempDirPrefix)
	if err != nil {
		log.Fatal(err)
	}
	if err = statefile.PutStatefiles(c.tmpDir, c.StateFiles); err != nil {
		return xerrors.Errorf("PutStatefiles err:%w", err)
	}
	return nil
}

// finalize remove temp dir
func (c *CheckPlugin) finalize() error {
	var err error
	c.StateFiles, err = statefile.GetStatefiles(c.tmpDir)
	if err != nil {
		return xerrors.Errorf("GetStatefiles err:%w", err)
	}
	err = c.PutCheckState(c.Rule.Name, c.CheckState)
	if err != nil {
		return xerrors.Errorf("PutCheckState err:%w", err)
	}
	return os.RemoveAll(c.tmpDir)

}

func (c *CheckPlugin) generate(ctx context.Context) (*mackerel.CheckReport, error) {
	cmd := cmdutil.CommandString(c.Rule.Command)
	now := time.Now()
	envs := append(c.Rule.Env, fmt.Sprintf("MACKEREL_PLUGIN_WORKDIR=%s", c.tmpDir))
	stdout, stderr, exitCode, err := cmdutil.RunCommand(ctx, cmd, "", envs, c.Rule.Timeout)

	if stderr != "" {
		log.Printf("plugin %s (%v): %q", c.Name, cmd, stderr)
	}

	var message string
	var status mackerel.CheckStatus
	if err != nil {
		log.Printf("Warning plugin %s (%v): %s", c.Name, cmd, err)
		message = err.Error()
		status = mackerel.CheckStatusUnknown
	} else {
		message = stdout
		status = exitCodeToStatus(exitCode)
	}

	newReport := mackerel.CheckReport{
		Source:               mackerel.NewCheckSourceHost(c.State.HostID),
		Name:                 c.Name,
		Status:               status,
		Message:              message,
		OccurredAt:           now.Unix(),
		NotificationInterval: c.Rule.NotificationInterval,
		MaxCheckAttempts:     c.Rule.MaxCheckAttempts,
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

func exitCodeToStatus(exitCode int) mackerel.CheckStatus {
	if code, ok := exitCodeMap[exitCode]; ok {
		return code
	}
	return mackerel.CheckStatusUnknown
}
