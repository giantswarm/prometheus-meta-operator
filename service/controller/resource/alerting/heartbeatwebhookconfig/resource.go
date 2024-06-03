package heartbeatwebhookconfig

import (
	"net/url"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1alpha1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1alpha1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
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
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.Client.MonitoringV1alpha1().AlertmanagerConfigs(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(v)
		},
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toAlertmanagerConfig(v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
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

func toAlertmanagerConfig(v interface{}, config Config) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := getObjectMeta(v)
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
	urlAddress, err := url.Parse(key.HeartbeatAPI(cluster, config.Installation))
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
	proxyURL, err := config.Proxy(opsgenieUrl)
	if err != nil {
		return nil, err
	}
	if proxyURL != nil {
		webhookHttpConfig.ProxyURL = proxyURL.String()
	}

	receiver := monitoringv1alpha1.Receiver{
		Name: key.HeartbeatReceiverName(cluster, config.Installation),
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
				Receiver: key.HeartbeatReceiverName(cluster, config.Installation),
				Matchers: []monitoringv1alpha1.Matcher{
					{Name: key.ClusterIDKey, Value: key.ClusterID(cluster)},
					{Name: key.InstallationKey, Value: config.Installation},
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

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*monitoringv1alpha1.AlertmanagerConfig)
	d := desired.(*monitoringv1alpha1.AlertmanagerConfig)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
