package legacy

import (
	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/controller"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	controllerresource "github.com/giantswarm/prometheus-meta-operator/service/controller/resource"
)

type ControllerConfig struct {
	Address                 string
	BaseDomain              string
	Bastions                []string
	Provider                string
	Installation            string
	Pipeline                string
	Region                  string
	Registry                string
	Customer                string
	CreatePVC               bool
	StorageSize             string
	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
	RetentionDuration       string
	RetentionSize           string
	OpsgenieKey             string
	RemoteWriteURL          string
	RemoteWriteUsername     string
	RemoteWritePassword     string
	K8sClient               k8sclient.Interface
	Logger                  micrologger.Logger
	PrometheusClient        promclient.Interface
	VpaClient               vpa_clientset.Interface
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
