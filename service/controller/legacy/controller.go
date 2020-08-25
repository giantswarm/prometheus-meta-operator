package legacy

import (
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/controller"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	controllerresource "github.com/giantswarm/prometheus-meta-operator/service/controller/resource"
)

type ControllerConfig struct {
	BaseDomain       string
	Provider         string
	CreatePVC        bool
	StorageSize      string
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
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
