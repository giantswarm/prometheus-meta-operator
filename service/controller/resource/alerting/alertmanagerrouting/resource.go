package alertmanagerrouting

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1alpha1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1alpha1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "alertmanagerrouting"
)

type Config struct {
	Client monitoringclient.Interface
	Logger micrologger.Logger

	Installation string
	HTTPProxy    string
	HTTPSProxy   string
	NoProxy      string
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
		GetObjectMeta: func(v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(v, config)
		},
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
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

func getObjectMeta(v interface{}, config Config) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.ClusterID(cluster),
		Namespace: key.Namespace(cluster),
		Labels:    key.AlertmanagerLabels(cluster),
	}, nil
}

func toAlertmanagerConfig(v interface{}, config Config) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := getObjectMeta(v, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	sendResolved := false
	urlAddress, err := url.Parse(fmt.Sprintf("https://api.opsgenie.com/v2/heartbeats/%s/ping", key.HeartbeatName(cluster, config.Installation)))
	if err != nil {
		return nil, microerror.Mask(err)
	}
	address := urlAddress.String()

	var proxyURL *string = nil
	if len(config.HTTPSProxy) > 0 {
		proxyURL = &config.HTTPSProxy
	} else if len(config.HTTPProxy) > 0 {
		proxyURL = &config.HTTPProxy
	}

	receiver := monitoringv1alpha1.Receiver{
		Name: key.HeartbeatReceiverName(cluster, config.Installation),
		WebhookConfigs: []monitoringv1alpha1.WebhookConfig{
			{
				URL: &address,
				HTTPConfig: &monitoringv1alpha1.HTTPConfig{
					ProxyURL: proxyURL,
					BasicAuth: &monitoringv1.BasicAuth{
						Password: corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "alertmanager-routing-secret",
							},
							Key: "opsgenie.key",
						},
					},
				},
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
					{Name: key.ClusterIDKey(), Value: key.ClusterID(cluster)},
					{Name: key.InstallationKey(), Value: config.Installation},
					{Name: key.TypeKey(), Value: key.Heartbeat()},
				},
				Continue: false,
				// We wait for 5 minutes before we start to ping OpsGenie to allow the prometheus server to start
				GroupWait:     "5m",
				GroupInterval: "30s",
				// We ping OpsGenie every minute
				RepeatInterval: "1m",
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
