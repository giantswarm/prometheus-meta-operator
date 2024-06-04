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

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/organization"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/template"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name                              = "scrapeconfigs"
	templateDirectory                 = "/opt/prometheus-meta-operator"
	templatePath                      = "files/templates/scrapeconfigs/*.yaml"
	unknownObservabilityBundleVersion = "0.0.0"
)

type Config struct {
	K8sClient          k8sclient.Interface
	Logger             micrologger.Logger
	OrganizationReader organization.Reader

	AdditionalScrapeConfigs   string
	Bastions                  []string
	Customer                  string
	Installation              string
	Pipeline                  string
	Provider                  cluster.Provider
	Region                    string
	Vault                     string
	WorkloadClusterETCDDomain string

	MimirEnabled bool
}

type Resource struct {
	k8sClient          k8sclient.Interface
	logger             micrologger.Logger
	organizationReader organization.Reader

	additionalScrapeConfigs   string
	bastions                  []string
	customer                  string
	installation              string
	pipeline                  string
	provider                  cluster.Provider
	region                    string
	vault                     string
	templatePath              string
	workloadClusterETCDDomain string

	mimirEnabled bool
}

type TemplateData struct {
	AdditionalScrapeConfigs   string
	APIServerURL              string
	AuthenticationType        string
	Bastions                  []string
	Pipeline                  string
	Provider                  string
	Region                    string
	ClusterID                 string
	ClusterType               string
	ServicePriority           string
	Customer                  string
	Organization              string
	SecretName                string
	EtcdSecretName            string
	Installation              string
	IgnoredTargets            string
	Vault                     string
	WorkloadClusterETCDDomain string
	CAPIManagementCluster     bool
	VintageManagementCluster  bool
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.OrganizationReader == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.OrganizationReader must not be empty")
	}

	return &Resource{
		k8sClient:          config.K8sClient,
		logger:             config.Logger,
		organizationReader: config.OrganizationReader,

		additionalScrapeConfigs:   config.AdditionalScrapeConfigs,
		bastions:                  config.Bastions,
		customer:                  config.Customer,
		installation:              config.Installation,
		pipeline:                  config.Pipeline,
		provider:                  config.Provider,
		region:                    config.Region,
		vault:                     config.Vault,
		workloadClusterETCDDomain: config.WorkloadClusterETCDDomain,

		templatePath: path.Join(templateDirectory, templatePath),

		mimirEnabled: config.MimirEnabled,
	}, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.PrometheusAdditionalScrapeConfigsSecretName(),
		Namespace: key.Namespace(cluster),
	}, nil
}

func (r *Resource) toSecret(ctx context.Context, v interface{}) (*corev1.Secret, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var workloadClusterETCDDomain string = ""
	if key.ClusterType(r.installation, v) == "workload_cluster" {
		clusterID := key.ClusterID(cluster)
		// Try to get the etcd url from the Giant Swarm way
		// TODO remove once all clusters are on CAPI
		service, err := r.k8sClient.K8sClient().CoreV1().Services(clusterID).Get(ctx, "master", metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// ETCD for CAPI clusters is monitored via service monitors
		} else if err != nil {
			return nil, microerror.Mask(err)
		} else {
			if value, ok := service.Annotations["giantswarm.io/etcd-domain"]; ok {
				workloadClusterETCDDomain = value
			}
		}
	}
	r.workloadClusterETCDDomain = workloadClusterETCDDomain

	scrapeConfigs, err := r.toData(ctx, v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	objectMeta, err := r.getObjectMeta(v)
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

func (r *Resource) toData(ctx context.Context, v interface{}) ([]byte, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	templateData, err := r.getTemplateData(ctx, cluster)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeConfigs, err := template.RenderTemplate(templateData, r.templatePath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return scrapeConfigs, nil
}

func (r *Resource) getTemplateData(ctx context.Context, cluster metav1.Object) (*TemplateData, error) {
	ignoredTargets, err := r.listTargetsToIgnore(ctx, cluster)

	if err != nil {
		return nil, microerror.Mask(err)
	}

	organization, err := r.organizationReader.Read(ctx, cluster)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	provider, err := key.ClusterProvider(cluster, r.provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var authenticationType = ""
	if !key.IsManagementCluster(r.installation, cluster) {
		authenticationType, err = key.ApiServerAuthenticationType(ctx, r.k8sClient, key.Namespace(cluster))
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	d := &TemplateData{
		AdditionalScrapeConfigs:   r.additionalScrapeConfigs,
		APIServerURL:              key.APIUrl(cluster),
		AuthenticationType:        authenticationType,
		Bastions:                  r.bastions,
		ClusterID:                 key.ClusterID(cluster),
		ClusterType:               key.ClusterType(r.installation, cluster),
		ServicePriority:           key.GetServicePriority(cluster),
		Customer:                  r.customer,
		Organization:              organization,
		Provider:                  provider,
		Pipeline:                  r.pipeline,
		Region:                    r.region,
		Installation:              r.installation,
		SecretName:                key.APIServerCertificatesSecretName,
		EtcdSecretName:            key.EtcdSecret(r.installation, cluster),
		Vault:                     r.vault,
		IgnoredTargets:            strings.Join(ignoredTargets[:], ","),
		WorkloadClusterETCDDomain: r.workloadClusterETCDDomain,
		CAPIManagementCluster:     key.IsCAPIManagementCluster(r.provider),
		VintageManagementCluster:  !key.IsCAPIManagementCluster(r.provider),
	}

	return d, nil
}

func (r *Resource) getObservabilityBundleAppVersion(ctx context.Context, cluster metav1.Object) (string, error) {
	appName := fmt.Sprintf("%s-observability-bundle", key.ClusterID(cluster))
	appNamespace := cluster.GetNamespace()

	if !key.IsCAPIManagementCluster(r.provider) {
		if key.IsManagementCluster(r.installation, cluster) {
			// Vintage MC
			appName = "observability-bundle"
			appNamespace = "giantswarm"
		} else {
			appNamespace = key.ClusterID(cluster)
		}
	}

	app := &appsv1alpha1.App{}
	objectKey := types.NamespacedName{Namespace: appNamespace, Name: appName}
	err := r.k8sClient.CtrlClient().Get(ctx, objectKey, app)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return unknownObservabilityBundleVersion, nil
		}
		return "", err
	}

	if app.Status.Version != "" {
		return app.Status.Version, nil
	}
	// This avoids a race condition where the app is created for the cluster but not deployed.
	return unknownObservabilityBundleVersion, nil
}

// List of targets we ignore in the scrape config (because they may be scraped by the agent or not scrappable)
func (r *Resource) listTargetsToIgnore(ctx context.Context, cluster metav1.Object) ([]string, error) {
	ignoredTargets := make([]string, 0)

	if key.IsEKSCluster(cluster) {
		// In case of EKS clusters, we assume scraping targets via ServiceMonitors,
		// so we ignore them from the Prometheus scrape config
		r.logger.Debugf(ctx, "EKS cluster: ignoring all scraping targets in Prometheus scrape config")
		ignoredTargets = append(ignoredTargets,
			"prometheus-operator-app",
			"kube-apiserver",
			"kube-controller-manager",
			"kube-scheduler",
			"node-exporter",
			"kubelet",
			"kube-proxy",
			"coredns",
			"kube-state-metrics",
			"etcd")
	} else {
		appVersion, err := r.getObservabilityBundleAppVersion(ctx, cluster)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		version, err := semver.Parse(appVersion)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		initialBundleVersion, err := semver.Parse("0.1.0")
		if err != nil {
			return nil, microerror.Mask(err)
		}

		bundleWithKSMAndExportersVersion, err := semver.Parse("0.4.0")
		if err != nil {
			return nil, microerror.Mask(err)
		}
		bundleWithKubeProxyExporterVersion, err := semver.Parse("0.8.3")
		if err != nil {
			return nil, microerror.Mask(err)
		}

		if version.GTE(initialBundleVersion) {
			ignoredTargets = append(ignoredTargets, "prometheus-operator-app", "kube-apiserver", "kube-controller-manager", "kube-scheduler", "node-exporter")
		}

		if version.GTE(bundleWithKSMAndExportersVersion) {
			ignoredTargets = append(ignoredTargets, "kubelet", "coredns", "kube-state-metrics")

			if key.IsCAPIManagementCluster(r.provider) {
				ignoredTargets = append(ignoredTargets, "etcd")
			}
		}
		if version.GTE(bundleWithKubeProxyExporterVersion) {
			ignoredTargets = append(ignoredTargets, "kube-proxy")
		}
	}
	// Vintage WC
	if !key.IsCAPIManagementCluster(r.provider) && !key.IsManagementCluster(r.installation, cluster) {
		// Since 18.0.0 we cannot scrape k8s endpoints externally so we ignore those targets.
		release := cluster.GetLabels()["release.giantswarm.io/version"]
		version, err := semver.Parse(release)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if version.Major >= 18 {
			ignoredTargets = append(ignoredTargets, "kube-controller-manager", "kube-scheduler")
		}
	}
	return ignoredTargets, nil
}

func (r *Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
