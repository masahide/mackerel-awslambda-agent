package config

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kelseyhightower/envconfig"
	mkconf "github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/masahide/mackerel-awslambda-agent/pkg/state"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store"
	"github.com/masahide/mackerel-awslambda-agent/pkg/store/dynamodbdriver"
	"github.com/pelletier/go-toml"
	"golang.org/x/xerrors"
)

const (
	maxMemo = 250
)

// AgentConfig is agent config struct.
type AgentConfig struct {
	*state.HostState
	Env
	CheckRules map[string]CheckRule
	*state.Manager
	hostStore  store.Store
	stateStore store.Store
	APIKey     string
}

// NewAgentConfig load config from env.
func NewAgentConfig(s store.Store, p client.ConfigProvider) (*AgentConfig, error) {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
	dynamo := dynamodbdriver.New(p, env.StateTable)
	a := &AgentConfig{
		stateStore: dynamo,
		hostStore:  dynamo,
		CheckRules: map[string]CheckRule{},
		Env:        env,
		Manager: &state.Manager{
			TTLDays:  env.StateTTLDays,
			Org:      env.Organization,
			Hostname: env.Hostname,
			Store:    s,
		},
	}
	return a, nil
}

func (a *AgentConfig) getCheckSum(param interface{}) string {
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("hosthash json.Marshal err:%s", err)
		return ""
	}
	h := fnv.New64a()
	// nolint:errcheck
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum64())
}

func (a *AgentConfig) GetHost() error {
	if a.HostState == nil {
		var err error
		a.HostState, err = a.GetHostState()
		if err != nil {
			return xerrors.Errorf("GetHostState err: %w", err)
		}
	}
	client := mackerel.NewClient(a.APIKey)
	checks := make([]mackerel.CheckConfig, 0, len(a.CheckRules))
	for _, rule := range a.CheckRules {
		checks = append(checks, mackerel.CheckConfig{Name: rule.Name, Memo: rule.Memo})
	}
	a.HostState.Hostname = a.Env.Hostname
	a.HostState.Organization = a.Env.Organization
	param := mackerel.CreateHostParam{
		Name:   a.Env.Hostname,
		Checks: checks,
	}
	if len(a.HostState.HostID) == 0 {
		id, err := client.CreateHost(&param)
		if err != nil {
			return xerrors.Errorf("mackerel CreateHost err: %w", err)
		}
		a.HostState.HostID = id
	}
	checkSum := a.getCheckSum(param)
	if checkSum == a.HostState.HostCheckSum {
		return nil
	}
	a.HostState.HostCheckSum = checkSum
	if err := a.PutHostState(*a.HostState); err != nil {
		return xerrors.Errorf("PutHostState err: %w", err)
	}
	return nil
}

func LoadS3Config(p client.ConfigProvider, s3Bucket, s3Key string) (*mkconf.Config, error) {
	svc := s3.New(p)
	xray.AWS(svc.Client)
	res, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &s3Bucket,
		Key:    &s3Key,
	})
	if err != nil {
		return nil, xerrors.Errorf("s3 GetObject bucket:%s key:%s err: %w", s3Bucket, s3Key, err)
	}
	defer res.Body.Close()
	conf := mkconf.Config{}
	err = toml.NewDecoder(res.Body).Decode(&conf)
	if err != nil {
		return nil, xerrors.Errorf("config toml decode err: %w", err)
	}
	return &conf, nil
}

func int32Value(v *int32) uint {
	if v == nil {
		return 0
	}
	return uint(*v)
}

func (a *AgentConfig) LoadAgentConfig(p client.ConfigProvider, s3Bucket, s3Key string) error {
	conf, err := LoadS3Config(p, s3Bucket, s3Key)
	if err != nil {
		return xerrors.Errorf("loadS3Config err: %w", err)
	}
	a.APIKey = conf.Apikey
	checkRules := map[string]CheckRule{}
	if pconfs, ok := conf.Plugin["checks"]; ok {
		for name, pconf := range pconfs {
			envs, err := pconf.Env.ConvertToStrings()
			if err != nil {
				envs = []string{}
			}
			rule := CheckRule{
				Name:                  name,
				Env:                   envs,
				Timeout:               time.Duration(pconf.TimeoutSeconds) * time.Second,
				CustomIdentifier:      aws.StringValue(pconf.CustomIdentifier),
				NotificationInterval:  int32Value(pconf.NotificationInterval.Minutes()),
				CheckInterval:         int(int32Value(pconf.CheckInterval.Minutes())),
				MaxCheckAttempts:      int32Value(pconf.MaxCheckAttempts),
				PreventAlertAutoClose: pconf.PreventAlertAutoClose,
			}
			switch v := pconf.CommandConfig.Raw.(type) {
			case string:
				rule.Command = v
			case []string:
				rule.Command = strings.Join(v, " ")
			default:
				log.Printf("Unknown command type rule:% value:%v", name, v)
			}
			if utf8.RuneCountInString(pconf.Memo) > maxMemo {
				log.Printf("plugin.checks.%s.memo' size exceeds 250 characters", name)
			}
			rule.Memo = pconf.Memo
			checkRules[name] = rule
		}
	}
	a.CheckRules = checkRules
	return nil
}
