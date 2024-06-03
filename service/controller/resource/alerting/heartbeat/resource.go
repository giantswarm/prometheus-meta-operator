package heartbeat

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/sirupsen/logrus"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "heartbeat"
)

type Config struct {
	Logger       micrologger.Logger
	Installation string
	OpsgenieKey  string
	Pipeline     string

	MimirEnabled bool
}

type Resource struct {
	logger          micrologger.Logger
	heartbeatClient *heartbeat.Client
	installation    string
	pipeline        string

	mimirEnabled bool
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}
	if config.Pipeline == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Pipeline must not be empty", config)
	}
	if config.OpsgenieKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.OpsgenieKey must not be empty", config)
	}

	c := &client.Config{
		ApiKey:         config.OpsgenieKey,
		OpsGenieAPIURL: client.API_URL,
		RetryCount:     1,
		LogLevel:       logrus.FatalLevel,
	}
	client, err := heartbeat.NewClient(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r := &Resource{
		logger:          config.Logger,
		heartbeatClient: client,
		installation:    config.Installation,
		pipeline:        config.Pipeline,
		mimirEnabled:    config.MimirEnabled,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toHeartbeat(v interface{}, installation string, pipeline string) (*heartbeat.Heartbeat, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// They need to be sorted alphabetically to avoid unnecessary heartbeat update
	tags := []string{
		"atlas",
		installation,
		fmt.Sprintf("managed-by: %s", project.Name()),
		fmt.Sprintf("pipeline: %s", pipeline),
	}
	sort.Strings(tags)

	h := &heartbeat.Heartbeat{
		Name:         key.HeartbeatName(cluster, installation),
		Description:  "📗 Runbook: https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/heartbeat-expired/",
		Interval:     60,
		IntervalUnit: string(heartbeat.Minutes),
		Enabled:      true,
		Expired:      false,
		OwnerTeam: og.OwnerTeam{
			Name: "alerts_router_team",
		},
		AlertTags:     tags,
		AlertPriority: "P3",
		AlertMessage:  fmt.Sprintf("Heartbeat [%s] is expired.", key.HeartbeatName(cluster, installation)),
	}

	return h, nil
}

func hasChanged(current, desired heartbeat.Heartbeat) bool {
	// Ignore those fields for comparison by setting them to the same value.
	current.Enabled = true
	desired.Enabled = true
	current.Expired = true
	desired.Expired = true
	// We get the ID back from opsgenie so we update it in the heartbeat
	desired.OwnerTeam.Id = current.OwnerTeam.Id

	return !reflect.DeepEqual(current, desired)
}
