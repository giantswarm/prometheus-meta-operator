package clusterapi

import (
	// If your operator watches a CRD import it here.
	// "github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"

	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/api/v1alpha2"

	controllerresource "github.com/giantswarm/prometheus-meta-operator/service/controller/resource"
)

type ControllerConfig struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface

	BaseDomain string
	Security   controllerresource.Security
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
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      "clusterapi-controller",
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(v1alpha2.Cluster)
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
