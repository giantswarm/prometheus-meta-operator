package alertmanagerconfig

import (
	"context"
	"io/ioutil"
	"path"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/domain"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/template"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name                     = "alertmanagerconfig"
	templateDirectory        = "/opt/prometheus-meta-operator"
	templatePath             = "files/templates/alertmanager/alertmanager.yaml"
	notificationTemplatePath = "files/templates/alertmanager/notification-template.tmpl"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	Installation       string
	Provider           string
	ProxyConfiguration domain.ProxyConfiguration
	OpsgenieKey        string
	GrafanaAddress     string
	SlackApiURL        string
	SlackProjectName   string
	Pipeline           string

	TemplatePath string
}

type TemplateData struct {
	Provider         string
	Installation     string
	ProxyURL         string
	OpsgenieKey      string
	GrafanaAddress   string
	SlackApiURL      string
	SlackProjectName string
	Pipeline         string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	if config.TemplatePath == "" {
		config.TemplatePath = path.Join(templateDirectory, templatePath)
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(v, config)
		},
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toSecret(v, config)
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
	return metav1.ObjectMeta{
		Name:      key.AlertManagerSecretName(),
		Namespace: key.NamespaceMonitoring(),
	}, nil
}

func toSecret(v interface{}, config Config) (*corev1.Secret, error) {
	objectMeta, err := getObjectMeta(v, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	notificationTemplate, err := ioutil.ReadFile(path.Join(templateDirectory, notificationTemplatePath))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	alertmanagerConfigSecret, err := toData(v, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			key.AlertmanagerKey():        alertmanagerConfigSecret,
			"notification-template.tmpl": notificationTemplate,
			key.OpsgenieKey():            []byte(config.OpsgenieKey),
		},
		Type: "Opaque",
	}

	return secret, nil
}

func toData(v interface{}, config Config) ([]byte, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	templateData, err := getTemplateData(cluster, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data, err := template.RenderTemplate(templateData, config.TemplatePath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return data, nil
}

func getTemplateData(cluster metav1.Object, config Config) (*TemplateData, error) {
	d := &TemplateData{
		Provider:         config.Provider,
		Installation:     config.Installation,
		ProxyURL:         config.ProxyConfiguration.GetURLForEndpoint("api.opsgenie.com"),
		OpsgenieKey:      config.OpsgenieKey,
		GrafanaAddress:   config.GrafanaAddress,
		SlackApiURL:      config.SlackApiURL,
		SlackProjectName: config.SlackProjectName,
		Pipeline:         config.Pipeline,
	}

	return d, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
