package scrapeconfigs

import (
	"bytes"
	"html/template"
	"path"
	"reflect"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name              = "scrapeconfigs"
	templateDirectory = "/opt/prometheus-meta-operator"
	templatePath      = "files/templates/*.yaml"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	Provider     string
	Vault        string
	TemplatePath string
}

type TemplateData struct {
	APIServerURL   string
	Provider       string
	ClusterID      string
	ClusterType    string
	SecretName     string
	EtcdSecretName string
	IsInCluster    bool
	Vault          string
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
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
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

func getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.PrometheusAdditionalScrapeConfigsSecretName(),
		Namespace: key.Namespace(cluster),
	}, nil
}

func toSecret(v interface{}, config Config) (*corev1.Secret, error) {
	scrapeConfigs, err := toData(v, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	objectMeta, err := getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeConfigsSecret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			key.PrometheusAdditionalScrapeConfigsName(): scrapeConfigs,
		},
		Type: "Opaque",
	}

	return scrapeConfigsSecret, nil
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

	scrapeConfigs, err := renderTemplate(*templateData, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return scrapeConfigs, nil
}

func getTemplateData(cluster metav1.Object, config Config) (*TemplateData, error) {
	clusterID := key.ClusterID(cluster)

	d := &TemplateData{
		APIServerURL:   key.APIUrl(cluster),
		ClusterID:      clusterID,
		ClusterType:    key.ClusterType(cluster),
		Provider:       config.Provider,
		SecretName:     key.Secret(),
		EtcdSecretName: key.EtcdSecret(cluster),
		IsInCluster:    key.IsInCluster(cluster),
		Vault:          config.Vault,
	}

	return d, nil
}

func renderTemplate(templateData TemplateData, config Config) ([]byte, error) {
	template, err := template.ParseGlob(config.TemplatePath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var b bytes.Buffer
	for _, t := range template.Templates() {
		err := t.Execute(&b, templateData)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	return b.Bytes(), nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
