package legacy

import (
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/resource"
	"k8s.io/apimachinery/pkg/runtime"

	controllerresource "github.com/giantswarm/prometheus-meta-operator/service/controller/resource"
)

type ControllerConfig struct {
	BaseDomain       string
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	Provider         string
}

type Controller struct {
	*controller.Controller
}

func NewController(config ControllerConfig) (*Controller, error) {
	var err error

	var resources []resource.Interface
	{
		c := controllerresource.Config{
			BaseDomain:       config.BaseDomain,
			K8sClient:        config.K8sClient,
			Logger:           config.Logger,
			PrometheusClient: config.PrometheusClient,
		}

		resources, err = controllerresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		var runtimeObject runtime.Object
		{
			switch config.Provider {
			case "aws":
				runtimeObject = new(v1alpha1.AWSConfig)
			case "azure":
				runtimeObject = new(v1alpha1.AzureConfig)
			case "kvm":
				runtimeObject = new(v1alpha1.KVMConfig)
			default:
				return nil, microerror.Maskf(invalidProviderError, "provider: %q", config.Provider)
			}
		}

		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      "legacy-controller",
			NewRuntimeObjectFunc: func() runtime.Object {
				return runtimeObject
			},
			Resources: resources,
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
