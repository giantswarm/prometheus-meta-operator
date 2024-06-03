package heartbeatwebhookconfig

import (
	"net/url"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1alpha1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1alpha1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "heartbeatwebhookconfig"
)

type Config struct {
	Client       monitoringclient.Interface
	Logger       micrologger.Logger
	Installation string
	Proxy        func(reqURL *url.URL) (*url.URL, error)

	MimirEnabled bool
}

type Resource struct {
	client       monitoringclient.Interface
	logger       micrologger.Logger
	installation string
	proxy        func(reqURL *url.URL) (*url.URL, error)

	mimirEnabled bool
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	r := &Resource{
		client:       config.Client,
		logger:       config.Logger,
		installation: config.Installation,
		proxy:        config.Proxy,
		mimirEnabled: config.MimirEnabled,
	}

	return r, nil
}

func (r Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.ClusterID(cluster),
		Namespace: key.MonitoringNamespace,
		Labels:    key.AlertmanagerLabels(),
	}, nil
}

func (r Resource) toAlertmanagerConfig(v interface{}) (*monitoringv1alpha1.AlertmanagerConfig, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	labels := make(map[string]string)
	for k, v := range key.AlertmanagerLabels() {
		labels[k] = v
	}

	sendResolved := false
	urlAddress, err := url.Parse(key.HeartbeatAPI(cluster, r.installation))
	if err != nil {
		return nil, microerror.Mask(err)
	}
	address := urlAddress.String()

	webhookHttpConfig := monitoringv1alpha1.HTTPConfig{
		Authorization: &monitoringv1.SafeAuthorization{
			Type: "GenieKey",
			Credentials: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: key.AlertmanagerGlobalSecretName,
				},
				Key: key.OpsGenieApiKey,
			},
		},
	}

	opsgenieUrl, err := url.Parse("https://api.opsgenie.com/v2/heartbeats")
	if err != nil {
		return nil, err
	}
	proxyURL, err := r.proxy(opsgenieUrl)
	if err != nil {
		return nil, err
	}
	if proxyURL != nil {
		webhookHttpConfig.ProxyURL = proxyURL.String()
	}

	receiver := monitoringv1alpha1.Receiver{
		Name: key.HeartbeatReceiverName(cluster, r.installation),
		WebhookConfigs: []monitoringv1alpha1.WebhookConfig{
			{
				URL:          &address,
				HTTPConfig:   &webhookHttpConfig,
				SendResolved: &sendResolved,
			},
		},
	}

	alertmanagerConfig := &monitoringv1alpha1.AlertmanagerConfig{
		ObjectMeta: objectMeta,
		Spec: monitoringv1alpha1.AlertmanagerConfigSpec{
			Route: &monitoringv1alpha1.Route{
				Receiver: key.HeartbeatReceiverName(cluster, r.installation),
				Matchers: []monitoringv1alpha1.Matcher{
					{Name: key.ClusterIDKey, Value: key.ClusterID(cluster)},
					{Name: key.InstallationKey, Value: r.installation},
					{Name: key.TypeKey, Value: key.Heartbeat()},
				},
				Continue: false,
				// wait for 30s before sending the first notification to opsgenie
				GroupWait: "30s",
				// wait for 30s between 2 alerts from the same group
				GroupInterval: "30s",
				// ping OpsGenie every 15 minutes
				RepeatInterval: "15m",
			},
			Receivers: []monitoringv1alpha1.Receiver{receiver},
		},
	}

	return alertmanagerConfig, nil
}

func (r Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*monitoringv1alpha1.AlertmanagerConfig)
	d := desired.(*monitoringv1alpha1.AlertmanagerConfig)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}

func (r *Resource) Name() string {
	return Name
}
