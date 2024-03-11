package clusterapi

import (
	"net/url"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"k8s.io/apimachinery/pkg/labels"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/client-go/dynamic"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
)

type ControllerConfig struct {
	K8sClient        k8sclient.Interface
	DynamicK8sClient dynamic.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	VpaClient        vpa_clientset.Interface
	Proxy            func(reqURL *url.URL) (*url.URL, error)

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	InsecureCA              bool
	Pipeline                string
	Provider                cluster.Provider
	Region                  string
	Registry                string

	OpsgenieKey string

	MimirEnabled bool

	PrometheusAddress            string
	PrometheusBaseDomain         string
	PrometheusEvaluationInterval string
	PrometheusLogLevel           string
	PrometheusScrapeInterval     string
	PrometheusImageRepository    string
	PrometheusVersion            string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string

	ExternalDNS bool
}

type Controller struct {
	*controller.Controller
}

func NewController(config ControllerConfig) (*Controller, error) {
	var err error

	var resources []resource.Interface
	{
		c := Config(config)

		resources, err = New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// Ensure cluster-api controller is not reconciling the MC cluster CR to avoid duplicate reconciliation with the management-cluster controller.
	// This due to CAPI installation providing a cluster CR for the MC.
	selector, err := labels.Parse("cluster.x-k8s.io/cluster-name!=" + config.Installation)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      project.Name() + "-cluster-api-controller",
			NewRuntimeObjectFunc: func() client.Object {
				return new(capi.Cluster)
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
