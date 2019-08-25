package plugin

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/kelseyhightower/envconfig"
	"github.com/masahide/mackerel-awslambda-agent/pkg/config"
	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
	"github.com/masahide/mackerel-awslambda-agent/pkg/statefile"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store/dynamodbdriver"
)

const (
	tempDirPrefix = "mackerel"
)

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
			Hostname: params.Hostname,
			Store:    dynamodbdriver.New(p, env.StateTable),
		},
	}
	return plugin
}

/*
func handler(ctx context.Context, params config.PluginParams) {
}
*/

// Init is load config of CheckPlugin
func (p *CheckPlugin) Init() error {
	var err error
	p.CheckState, err = p.GetCheckState(p.CheckRule.Name)
	if err != nil {
		return err
	}
	p.tmpDir, err = ioutil.TempDir(p.Env.TempDir, tempDirPrefix)
	if err != nil {
		log.Fatal(err)
	}
	if err = statefile.PutStatefiles(p.tmpDir, p.StateFiles); err != nil {
		return err
	}
	return nil
}

// Cleanup remove tmpdir
func (p *CheckPlugin) Cleanup() error {
	var err error
	p.StateFiles, err = statefile.GetStatefiles(p.tmpDir)
	if err != nil {
		return err
	}
	err = p.PutCheckState(p.CheckRule.Name, p.CheckState)
	if err != nil {
		return err
	}
	return os.RemoveAll(p.tmpDir)

}
