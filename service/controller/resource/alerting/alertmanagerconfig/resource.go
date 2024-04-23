package alertmanagerconfig

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/template"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name                     = "alertmanagerconfig"
	templateDirectory        = "/opt/prometheus-meta-operator"
	alertmanagerTemplatePath = "files/templates/alertmanager/alertmanager.yaml"
	notificationTemplatePath = "files/templates/alertmanager/notification-template.tmpl"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	BaseDomain     string
	GrafanaAddress string
	Installation   string
	MimirEnabled   bool
	OpsgenieKey    string
	Pipeline       string
	Proxy          func(reqURL *url.URL) (*url.URL, error)
	SlackApiToken  string
}

type NotificationTemplateData struct {
	GrafanaAddress    string
	MimirEnabled      bool
	PrometheusAddress string
}

type AlertmanagerTemplateData struct {
	Installation  string
	OpsgenieKey   string
	Pipeline      string
	ProxyURL      string
	SlackApiToken string
	MimirEnabled  bool
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta()
		},
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toSecret(config)
		},
		HasChangedFunc: hasChanged,
	}

	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta() (metav1.ObjectMeta, error) {
	return metav1.ObjectMeta{
		Name:      key.AlertmanagerSecretName(),
		Namespace: key.MonitoringNamespace,
	}, nil
}

func toSecret(config Config) (*corev1.Secret, error) {
	objectMeta, err := getObjectMeta()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	notificationTemplate, err := renderNotificationTemplate(templateDirectory, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	alertmanagerConfigSecret, err := renderAlertmanagerConfig(templateDirectory, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			"alertmanager.yaml":          alertmanagerConfigSecret,
			"notification-template.tmpl": notificationTemplate,
		},
		Type: "Opaque",
	}

	return secret, nil
}

func renderNotificationTemplate(templateDirectory string, config Config) ([]byte, error) {
	templateData := NotificationTemplateData{
		GrafanaAddress:    config.GrafanaAddress,
		MimirEnabled:      config.MimirEnabled,
		PrometheusAddress: fmt.Sprintf("https://%s", config.BaseDomain),
	}

	data, err := template.RenderTemplate(templateData, path.Join(templateDirectory, notificationTemplatePath))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return data, nil
}

func renderAlertmanagerConfig(templateDirectory string, config Config) ([]byte, error) {
	templateData, err := getTemplateData(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data, err := template.RenderTemplate(templateData, path.Join(templateDirectory, alertmanagerTemplatePath))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return data, nil
}

func getTemplateData(config Config) (*AlertmanagerTemplateData, error) {
	opsgenieUrl, err := url.Parse("https://api.opsgenie.com/v2/heartbeats")
	if err != nil {
		return nil, err
	}
	proxyURL, err := config.Proxy(opsgenieUrl)
	if err != nil {
		return nil, err
	}

	d := &AlertmanagerTemplateData{
		Installation:  config.Installation,
		OpsgenieKey:   config.OpsgenieKey,
		Pipeline:      config.Pipeline,
		SlackApiToken: config.SlackApiToken,
		MimirEnabled:  config.MimirEnabled,
	}

	if proxyURL != nil {
		d.ProxyURL = proxyURL.String()
	}

	return d, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
