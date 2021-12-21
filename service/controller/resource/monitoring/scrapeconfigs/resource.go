package scrapeconfigs

import (
	"context"
	"path"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/pkg/template"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name              = "scrapeconfigs"
	templateDirectory = "/opt/prometheus-meta-operator"
	templatePath      = "files/templates/scrapeconfigs/*.yaml"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	AdditionalScrapeConfigs   string
	Bastions                  []string
	Installation              string
	Provider                  string
	Mayu                      string
	Vault                     string
	TemplatePath              string
	WorkloadClusterETCDDomain string
}

type TemplateData struct {
	AdditionalScrapeConfigs   string
	APIServerURL              string
	Bastions                  []string
	Provider                  string
	ClusterID                 string
	ClusterType               string
	SecretName                string
	EtcdSecretName            string
	Installation              string
	IsInCluster               bool
	Mayu                      string
	Vault                     string
	WorkloadClusterETCDDomain string
	CAPICluster               bool
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
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var workloadClusterETCDDomain string = ""
	if "workload_cluster" == key.ClusterType(v) {
		clusterID := key.ClusterID(cluster)
		// Try to get the etcd url from the Giant Swarm way
		service, err := config.K8sClient.K8sClient().CoreV1().Services(clusterID).Get(context.Background(), "master", metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// TODO we ignore ETCD for capi clusters for now. Find a way to do it later
		} else if err != nil {
			return nil, microerror.Mask(err)
		} else {
			if value, ok := service.Annotations["giantswarm.io/etcd-domain"]; ok {
				workloadClusterETCDDomain = value
			}
		}
	}
	config.WorkloadClusterETCDDomain = workloadClusterETCDDomain

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

	scrapeConfigs, err := template.RenderTemplate(templateData, config.TemplatePath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return scrapeConfigs, nil
}

func getTemplateData(cluster metav1.Object, config Config) (*TemplateData, error) {
	clusterID := key.ClusterID(cluster)

	d := &TemplateData{
		AdditionalScrapeConfigs:   config.AdditionalScrapeConfigs,
		APIServerURL:              key.APIUrl(cluster),
		Bastions:                  config.Bastions,
		ClusterID:                 clusterID,
		ClusterType:               key.ClusterType(cluster),
		Provider:                  config.Provider,
		Installation:              config.Installation,
		SecretName:                key.Secret(),
		EtcdSecretName:            key.EtcdSecret(cluster),
		Vault:                     config.Vault,
		Mayu:                      config.Mayu,
		WorkloadClusterETCDDomain: config.WorkloadClusterETCDDomain,
		CAPICluster:               key.IsCAPICluster(cluster),
	}

	return d, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
