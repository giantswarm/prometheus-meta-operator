package legacy

import (
	"github.com/giantswarm/apiextensions/v3/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v4/pkg/controller"
	"github.com/giantswarm/operatorkit/v4/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

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

	Bastions     []string
	Customer     string
	Installation string
	Pipeline     string
	Provider     string
	Region       string
	Registry     string

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

	var operatorkitController *controller.Controller
	{
		var runtimeFunc func() runtime.Object
		{
			switch config.Provider {
			case "aws":
				runtimeFunc = func() runtime.Object {
					return new(v1alpha1.AWSConfig)
				}
			case "azure":
				runtimeFunc = func() runtime.Object {
					return new(v1alpha1.AzureConfig)
				}
			case "kvm":
				runtimeFunc = func() runtime.Object {
					return new(v1alpha1.KVMConfig)
				}
			default:
				return nil, microerror.Maskf(invalidProviderError, "provider: %q", config.Provider)
			}
		}

		c := controller.Config{
			K8sClient:            config.K8sClient,
			Logger:               config.Logger,
			Name:                 project.Name() + "-legacy-controller",
			NewRuntimeObjectFunc: runtimeFunc,
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
