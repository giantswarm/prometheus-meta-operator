package heartbeat

import (
	"fmt"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "heartbeat"
)

type Config struct {
	Logger       micrologger.Logger
	Installation string
}

type Resource struct {
	logger          micrologger.Logger
	heartbeatClient *heartbeat.Client
	installation    string
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	c := &client.Config{
		ApiKey:         "",
		OpsGenieAPIURL: client.API_URL,
		RetryCount:     1,
	}
	client, err := heartbeat.NewClient(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r := &Resource{
		logger:          config.Logger,
		heartbeatClient: client,
		installation:    config.Installation,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toHeartbeat(v interface{}, installation string) (*heartbeat.Heartbeat, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	name := fmt.Sprintf("%s-%s", installation, cluster.GetName())
	h := &heartbeat.Heartbeat{
		Name:         name,
		Description:  "*Recipe:* https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/heartbeat-expired/",
		Interval:     25,
		IntervalUnit: string(heartbeat.Minutes),
		Enabled:      true,
		Expired:      false,
		OwnerTeam: og.OwnerTeam{
			Name: "alerts_router_team",
		},
		AlertTags: []string{
			"managed-by: prometheus-meta-operator",
		},
		AlertPriority: "P3",
		AlertMessage:  fmt.Sprintf("Heartbeat [%s] is expired.", name),
	}

	return h, nil
}

func hasChanged(current, desired heartbeat.Heartbeat) bool {
	current.Enabled = true
	desired.Enabled = true
	current.Expired = true
	desired.Expired = true

	return !reflect.DeepEqual(current, desired)
}
