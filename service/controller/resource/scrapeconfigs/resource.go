package scrapeconfigs

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/api/v1alpha2"

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
	ETCD         string
	Provider     string
	ClusterID    string
	SecretName   string
	IsInCluster  bool
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
			return toSecret(v, config.Provider, config.K8sClient)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toSecret(v interface{}, provider string, clients k8sclient.Interface) (*corev1.Secret, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	templateData, err := getTemplateData(cluster, provider, clients)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeConfigs, err := renderTemplate(*templateData)
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

func getTemplateData(cluster metav1.Object, provider string, clients k8sclient.Interface) (*TemplateData, error) {
	var etcd string
	switch v := cluster.(type) {
	case *v1alpha2.Cluster:
		ctx := context.Background()
		infra, err := clients.G8sClient().InfrastructureV1alpha2().AWSClusters(v.Spec.InfrastructureRef.Namespace).Get(ctx, v.Spec.InfrastructureRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, microerror.Mask(err)
		}

		etcd = fmt.Sprintf("etcd.%s.k8s.%s:2379", key.ClusterID(cluster), infra.Spec.Cluster.DNS.Domain)
	case *v1alpha1.AWSConfig:
		etcd = fmt.Sprintf("%s:%d", v.Spec.Cluster.Etcd.Domain, v.Spec.Cluster.Etcd.Port)
	case *v1alpha1.AzureConfig:
		etcd = fmt.Sprintf("%s:%d", v.Spec.Cluster.Etcd.Domain, v.Spec.Cluster.Etcd.Port)
	case *v1alpha1.KVMConfig:
		etcd = fmt.Sprintf("%s:%d", v.Spec.Cluster.Etcd.Domain, v.Spec.Cluster.Etcd.Port)
	default:
		return nil, microerror.Maskf(wrongTypeError, fmt.Sprintf("%T", v))
	}

	clusterID := key.ClusterID(cluster)

	d := &TemplateData{
		APIServerURL: key.APIUrl(cluster),
		ClusterID:    clusterID,
		Provider:     provider,
		SecretName:   key.Secret(),
		ETCD:         etcd,
		IsInCluster:  key.IsInCluster(cluster),
	}

	return d, nil
}

func renderTemplate(templateData TemplateData) (string, error) {
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
