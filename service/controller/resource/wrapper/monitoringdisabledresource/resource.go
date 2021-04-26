package monitoringdisabledresource

import (
	"context"
	"net/url"
	"path"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v4/pkg/resource"
	clientruntime "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/silence"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	Resource resource.Interface
	Logger   micrologger.Logger
}

type monitoringDisabledWrapper struct {
	resource resource.Interface
	logger   micrologger.Logger
}

// New returns a new monitoring disabled wrapper according to the configured resource's
// implementation, which might be resource.Interface or crud.Interface. This has
// then different implications on how to measure metrics for the different
// methods of the interfaces.
func New(config Config) (resource.Interface, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &monitoringDisabledWrapper{
		resource: config.Resource,
		logger:   config.Logger,
	}

	return r, nil
}

func (r *monitoringDisabledWrapper) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	var isDisabled = false

	if key.IsMonitoringDisabled(cluster) {
		r.logger.Debugf(ctx, "monitoring disabled, cleaning up existing monitoring")
		isDisabled = true
	}

	{
		alertmanagerURL, _ := url.Parse("")
		silenceParams := silence.NewGetSilencesParams().WithContext(ctx)
		amclient := NewAlertmanagerClient(alertmanagerURL)

		getOk, err := amclient.Silence.GetSilences(silenceParams)
		if err != nil {
			return err
		}

		for _, silence := range getOk.Payload {
			// TODO: find whole silence / silence for entire cluster
			if silence.ID != nil {
			}
		}
	}

	if isDisabled {
		return r.resource.EnsureDeleted(ctx, obj)
	}

	return r.resource.EnsureCreated(ctx, obj)
}

const (
	defaultAmHost      = "localhost"
	defaultAmPort      = "9093"
	defaultAmApiv2path = "/api/v2"
)

// NewAlertmanagerClient initializes an alertmanager client with the given URL
func NewAlertmanagerClient(amURL *url.URL) *client.Alertmanager {
	address := defaultAmHost + ":" + defaultAmPort
	schemes := []string{"http"}

	if amURL.Host != "" {
		address = amURL.Host // URL documents host as host or host:port
	}
	if amURL.Scheme != "" {
		schemes = []string{amURL.Scheme}
	}

	cr := clientruntime.New(address, path.Join(amURL.Path, defaultAmApiv2path), schemes)

	if amURL.User != nil {
		password, _ := amURL.User.Password()
		cr.DefaultAuthentication = clientruntime.BasicAuth(amURL.User.Username(), password)
	}

	return client.New(cr, strfmt.Default)
}

func (r *monitoringDisabledWrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return r.resource.EnsureDeleted(ctx, obj)
}

func (r *monitoringDisabledWrapper) Name() string {
	return r.resource.Name()
}
