package controlplane

import (
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/controller"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
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
	PrometheusVersion       string
	Customer                string
	CreatePVC               bool
	StorageSize             string
	Vault                   string
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
		c := resourcesConfig(config)

		resources, err = newResources(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var selector controller.Selector
	{
		set := labels.Set{
			"component": "apiserver",
		}
		selector = set.AsSelector()
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      project.Name() + "-control-plane-controller",
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(v1.Service)
			},
			Resources: resources,
			Selector:  selector,
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
