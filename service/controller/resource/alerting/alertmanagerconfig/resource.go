package alertmanagerconfig

import (
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
	SlackApiURL    string
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	baseDomain     string
	grafanaAddress string
	installation   string
	mimirEnabled   bool
	opsgenieKey    string
	pipeline       string
	proxy          func(reqURL *url.URL) (*url.URL, error)
	slackApiToken  string
	slackApiURL    string
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
	SlackApiURL   string
	MimirEnabled  bool
}

func New(config Config) (*Resource, error) {
	return &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		baseDomain:     config.BaseDomain,
		grafanaAddress: config.GrafanaAddress,
		installation:   config.Installation,
		mimirEnabled:   config.MimirEnabled,
		opsgenieKey:    config.OpsgenieKey,
		pipeline:       config.Pipeline,
		proxy:          config.Proxy,
		slackApiToken:  config.SlackApiToken,
		slackApiURL:    config.SlackApiURL,
	}, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) getObjectMeta() (metav1.ObjectMeta, error) {
	return metav1.ObjectMeta{
		Name:      key.AlertmanagerSecretName(),
		Namespace: key.MonitoringNamespace,
	}, nil
}

func (r *Resource) toSecret() (*corev1.Secret, error) {
	objectMeta, err := r.getObjectMeta()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	notificationTemplate, err := r.renderNotificationTemplate(templateDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	alertmanagerConfigSecret, err := r.renderAlertmanagerConfig(templateDirectory)
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

func (r *Resource) renderNotificationTemplate(templateDirectory string) ([]byte, error) {
	templateData := NotificationTemplateData{
		GrafanaAddress:    r.grafanaAddress,
		MimirEnabled:      r.mimirEnabled,
		PrometheusAddress: fmt.Sprintf("https://%s", r.baseDomain),
	}

	data, err := template.RenderTemplate(templateData, path.Join(templateDirectory, notificationTemplatePath))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return data, nil
}

func (r *Resource) renderAlertmanagerConfig(templateDirectory string) ([]byte, error) {
	templateData, err := r.getTemplateData()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data, err := template.RenderTemplate(templateData, path.Join(templateDirectory, alertmanagerTemplatePath))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return data, nil
}

func (r *Resource) getTemplateData() (*AlertmanagerTemplateData, error) {
	opsgenieUrl, err := url.Parse("https://api.opsgenie.com/v2/heartbeats")
	if err != nil {
		return nil, err
	}
	proxyURL, err := r.proxy(opsgenieUrl)
	if err != nil {
		return nil, err
	}

	d := &AlertmanagerTemplateData{
		Installation:  r.installation,
		OpsgenieKey:   r.opsgenieKey,
		Pipeline:      r.pipeline,
		SlackApiToken: r.slackApiToken,
		SlackApiURL:   r.slackApiURL,
		MimirEnabled:  r.mimirEnabled,
	}

	if proxyURL != nil {
		d.ProxyURL = proxyURL.String()
	}

	return d, nil
}

func (r *Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
