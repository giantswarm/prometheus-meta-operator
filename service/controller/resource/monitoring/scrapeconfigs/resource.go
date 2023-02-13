package scrapeconfigs

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/blang/semver"
	appsv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/template"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name              = "scrapeconfigs"
	templateDirectory = "/opt/prometheus-meta-operator"
	templatePath      = "files/templates/scrapeconfigs/*.yaml"
)

var kubernetesTargets = []string{"kube-apiserver", "kube-controller-manager", "kube-scheduler"}

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	AdditionalScrapeConfigs   string
	Bastions                  []string
	Customer                  string
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
	ServicePriority           string
	Customer                  string
	Organization              string
	SecretName                string
	EtcdSecretName            string
	Installation              string
	IgnoredTargets            string
	Mayu                      string
	Vault                     string
	WorkloadClusterETCDDomain string
	CAPICluster               bool
	CAPIManagementCluster     bool
	VintageManagementCluster  bool
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
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toSecret(ctx, v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.PrometheusAdditionalScrapeConfigsSecretName(),
		Namespace: key.Namespace(cluster),
	}, nil
}

func toSecret(ctx context.Context, v interface{}, config Config) (*corev1.Secret, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var workloadClusterETCDDomain string = ""
	if key.ClusterType(config.Installation, v) == "workload_cluster" {
		clusterID := key.ClusterID(cluster)
		// Try to get the etcd url from the Giant Swarm way
		service, err := config.K8sClient.K8sClient().CoreV1().Services(clusterID).Get(ctx, "master", metav1.GetOptions{})
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

	scrapeConfigs, err := toData(ctx, config.K8sClient.CtrlClient(), v, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	objectMeta, err := getObjectMeta(ctx, v)
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

func toData(ctx context.Context, ctrlClient client.Client, v interface{}, config Config) ([]byte, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	templateData, err := getTemplateData(ctx, ctrlClient, cluster, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeConfigs, err := template.RenderTemplate(templateData, config.TemplatePath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return scrapeConfigs, nil
}

func getTemplateData(ctx context.Context, ctrlClient client.Client, cluster metav1.Object, config Config) (*TemplateData, error) {
	ignoredTargets, err := listTargetsToIgnore(ctx, ctrlClient, cluster, config)

	if err != nil {
		return nil, microerror.Mask(err)
	}

	d := &TemplateData{
		AdditionalScrapeConfigs:   config.AdditionalScrapeConfigs,
		APIServerURL:              key.APIUrl(cluster),
		Bastions:                  config.Bastions,
		ClusterID:                 key.ClusterID(cluster),
		ClusterType:               key.ClusterType(config.Installation, cluster),
		ServicePriority:           key.GetServicePriority(cluster),
		Customer:                  config.Customer,
		Organization:              key.GetOrganization(cluster),
		Provider:                  config.Provider,
		Installation:              config.Installation,
		SecretName:                key.Secret(),
		EtcdSecretName:            key.EtcdSecret(config.Installation, cluster),
		Vault:                     config.Vault,
		Mayu:                      config.Mayu,
		IgnoredTargets:            strings.Join(ignoredTargets[:], ","),
		WorkloadClusterETCDDomain: config.WorkloadClusterETCDDomain,
		CAPICluster:               key.IsCAPICluster(cluster),
		CAPIManagementCluster:     key.IsCAPIManagementCluster(config.Provider),
		VintageManagementCluster:  !key.IsCAPIManagementCluster(config.Provider),
	}

	return d, nil
}

func getDefaultAppVersion(ctx context.Context, ctrlClient client.Client, cluster metav1.Object, config Config) (string, error) {
	appName := fmt.Sprintf("%s-default-apps", key.ClusterID(cluster))
	appNamespace := cluster.GetNamespace()
	objectKey := types.NamespacedName{Namespace: appNamespace, Name: appName}

	app := &appsv1alpha1.App{}
	err := ctrlClient.Get(ctx, objectKey, app)
	if err != nil {
		return "", err
	}
	return app.Status.Version, nil
}

// List of targets we ignore in the scrape config (whether they are scraped by the agent or they are unscrappable)
func listTargetsToIgnore(ctx context.Context, ctrlClient client.Client, cluster metav1.Object, config Config) ([]string, error) {
	ignoredTargets := make([]string, 0)

	// For CAPI clusters, this is a case to case basis. We need to check the default app version for now.
	if key.IsCAPIManagementCluster(config.Provider) {
		ignoredTargets = append(ignoredTargets, "docker", "vault")
		appVersion, err := getDefaultAppVersion(ctx, ctrlClient, cluster, config)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		version, err := semver.Parse(appVersion)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		switch config.Provider {
		case "capa":
			capaAgentVersion, err := semver.Parse("0.11.0")
			if err != nil {
				return nil, microerror.Mask(err)
			}
			if version.GTE(capaAgentVersion) {
				ignoredTargets = append(ignoredTargets, kubernetesTargets...)
			}
		case "cloud-directory":
			cloudDirectorAgentVersion, err := semver.Parse("0.3.0")
			if err != nil {
				return nil, microerror.Mask(err)
			}
			if version.GTE(cloudDirectorAgentVersion) {
				ignoredTargets = append(ignoredTargets, kubernetesTargets...)
			}
		case "gcp":
			gcpAgentVersion, err := semver.Parse("0.16.0")
			if err != nil {
				return nil, microerror.Mask(err)
			}
			if version.GTE(gcpAgentVersion) {
				ignoredTargets = append(ignoredTargets, kubernetesTargets...)
			}
		case "openstack":
			openstackAgentVersion, err := semver.Parse("0.8.0")
			if err != nil {
				return nil, microerror.Mask(err)
			}
			if version.GTE(openstackAgentVersion) {
				ignoredTargets = append(ignoredTargets, kubernetesTargets...)
			}
		}
	} else if key.IsManagementCluster(config.Installation, cluster) { // Vintage MC
		// Vintage MCs (apart from KVM) are all running agents.
		ignoredTargets = append(ignoredTargets, kubernetesTargets...)
	} else { // Vintage WC
		// Since 18.0.0 we cannot scrape k8s endpoints externally so we ignore those targets.
		release := cluster.GetLabels()["release.giantswarm.io/version"]
		version, err := semver.Parse(release)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if version.Major >= 19 ||
			(config.Provider == "aws" && version.Major == 18 && version.Minor >= 2) { // aws 18.2.0 is running the agent
			ignoredTargets = append(ignoredTargets, kubernetesTargets...)
		} else if version.Major >= 18 {
			ignoredTargets = append(ignoredTargets, "kube-controller-manager", "kube-scheduler")
		}
	}
	return ignoredTargets, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
