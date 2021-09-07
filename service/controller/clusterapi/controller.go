package clusterapi

import (
	"context"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v4/pkg/controller"
	"github.com/giantswarm/operatorkit/v4/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	capiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	controllerresource "github.com/giantswarm/prometheus-meta-operator/service/controller/resource"
)

type ControllerConfig struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	VpaClient        vpa_clientset.Interface

	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	Pipeline                string
	Provider                string
	Region                  string
	Registry                string

	OpsgenieKey string

	PrometheusAddress             string
	PrometheusBaseDomain          string
	PrometheusCreatePVC           bool
	PrometheusStorageSize         string
	PrometheusLogLevel            string
	PrometheusRemoteWriteURL      string
	PrometheusRemoteWriteUsername string
	PrometheusRemoteWritePassword string
	PrometheusRetentionDuration   string
	PrometheusRetentionSize       string
	PrometheusVersion             string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
}

type Controller struct {
	*controller.Controller
}

func NewController(config ControllerConfig) (*Controller, error) {
	var err error

	var resources []resource.Interface
	{
		c := controllerresource.Config(config)

		resources, err = controllerresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	runtimeObjectFactoryFunc, err := getClusterFactoryFunc(config.K8sClient.CtrlClient())
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient:            config.K8sClient,
			Logger:               config.Logger,
			Name:                 project.Name() + "-cluster-api-controller",
			NewRuntimeObjectFunc: runtimeObjectFactoryFunc,
			Resources:            resources,
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Controller{
		Controller: operatorkitController,
	}

	return c, nil
}

func getClusterFactoryFunc(ctrlClient client.Client) (func() runtime.Object, error) {
	var clusterCRD apiextensionsv1.CustomResourceDefinition
	err := ctrlClient.Get(context.Background(), client.ObjectKey{
		Name: "clusters.cluster.x-k8s.io",
	}, &clusterCRD)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Find out configured storage version.
	var storageVersion string
	for _, v := range clusterCRD.Spec.Versions {
		if v.Storage {
			storageVersion = v.Name
			break
		}
	}

	// Decide which object to construct based on storage version.
	var fn func() runtime.Object
	switch storageVersion {
	case "v1alpha3":
		fn = func() runtime.Object { return new(capiv1alpha3.Cluster) }
	case "v1alpha4":
		fn = func() runtime.Object { return new(capiv1alpha4.Cluster) }
	default:
		return nil, microerror.Maskf(unsupportedStorageVersionError, "implementation does not support storage version %q", storageVersion)
	}

	return fn, nil
}
