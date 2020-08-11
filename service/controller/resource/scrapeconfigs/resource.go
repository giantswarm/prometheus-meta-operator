package scrapeconfigs

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/pkg/templates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name         = "scrapeconfigs"
	templatePath = "/opt/prometheus-meta-operator/files/templates/additional-scrape-configs.template.yaml"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	Provider  string
}

type TemplateData struct {
	APIServerURL string
	Provider     string
	ClusterID    string
	SecretName   string
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
		ToCR: func(v interface{}) (metav1.Object, error) {
			return toSecret(v, config.Provider)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toSecret(v interface{}, provider string) (*corev1.Secret, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var clusterID = cluster.GetName()
	scrapeConfigs, err := renderTemplate(clusterID, provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeConfigsSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.PrometheusAdditionalScrapeConfigsSecretName(),
			Namespace: key.Namespace(cluster),
		},
		Data: map[string][]byte{
			key.PrometheusAdditionalScrapeConfigsName(): []byte(scrapeConfigs),
		},
		Type: "Opaque",
	}

	return scrapeConfigsSecret, nil
}

func renderTemplate(clusterID string, provider string) (string, error) {
	var templateData = TemplateData{
		APIServerURL: fmt.Sprintf("master.%s", clusterID),
		ClusterID:    clusterID,
		Provider:     provider,
		SecretName:   key.Secret(),
	}

	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", microerror.Mask(err)
	}

	template := string(content)

	scrapeConfigs, err := templates.Render(template, templateData)
	if err != nil {
		return "", microerror.Mask(err)
	}
	return scrapeConfigs, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
