package alertmanagerconfig

import (
	_ "embed"
	"fmt"
	htmltemplate "html/template"
	"net/url"
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

var (
	//go:embed templates/alertmanager.yaml
	alertmanagerConfig         string
	alertmanagerConfigTemplate *htmltemplate.Template

	//go:embed templates/notification-template.tmpl
	notificationTemplate         string
	notificationTemplateTemplate *htmltemplate.Template
)

func init() {
	alertmanagerConfigTemplate = htmltemplate.Must(template.New("alertmanager.yaml").Parse(alertmanagerConfig))
	notificationTemplateTemplate = htmltemplate.Must(template.New("notification-template.tmpl").Parse(notificationTemplate))
}

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	AlertmanagerEnabled bool

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
	config Config
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
	return &Resource{config}, nil
}

func (r *Resource) Name() string {
	return Name
}

func getObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      key.AlertmanagerSecretName(),
		Namespace: key.MonitoringNamespace,
	}
}

func (r *Resource) toSecret() (*corev1.Secret, error) {
	notificationTemplate, err := r.renderNotificationTemplate(templateDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	alertmanagerConfigSecret, err := r.renderAlertmanagerConfig(templateDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: getObjectMeta(),
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
		GrafanaAddress:    r.config.GrafanaAddress,
		MimirEnabled:      r.config.MimirEnabled,
		PrometheusAddress: fmt.Sprintf("https://%s", r.config.BaseDomain),
	}

	data, err := template.Render(notificationTemplateTemplate, templateData)
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

	data, err := template.Render(alertmanagerConfigTemplate, templateData)
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
	proxyURL, err := r.config.Proxy(opsgenieUrl)
	if err != nil {
		return nil, err
	}

	d := &AlertmanagerTemplateData{
		Installation:  r.config.Installation,
		OpsgenieKey:   r.config.OpsgenieKey,
		Pipeline:      r.config.Pipeline,
		SlackApiToken: r.config.SlackApiToken,
		SlackApiURL:   r.config.SlackApiURL,
		MimirEnabled:  r.config.MimirEnabled,
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
