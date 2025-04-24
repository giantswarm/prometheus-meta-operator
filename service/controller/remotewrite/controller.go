package remotewrite

import (
	"net/url"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"

	"sigs.k8s.io/controller-runtime/pkg/client"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
)

type ControllerConfig struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	Proxy            func(reqURL *url.URL) (*url.URL, error)
}

type Controller struct {
	*controller.Controller
}

func NewController(config ControllerConfig) (*Controller, error) {
	var err error

	var resources []resource.Interface
	{
		resources, err = newResources(config)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      project.Name() + "-remote-write-controller",
			NewRuntimeObjectFunc: func() client.Object {
				return new(pmov1alpha1.RemoteWrite)
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
