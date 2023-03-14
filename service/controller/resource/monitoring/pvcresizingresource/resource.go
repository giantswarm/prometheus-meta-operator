package pvcresizingresource

import (
	"fmt"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/micrologger"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
)

const (
	Name = "pvcresizingresource"
)

type Config struct {
	K8sClient        k8sclient.Interface
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
}

type Resource struct {
	k8sClient  k8sclient.Interface
	promClient promclient.Interface
	logger     micrologger.Logger
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient:  config.K8sClient,
		promClient: config.PrometheusClient,
		logger:     config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func resourceName(clusterID string) string {
	return fmt.Sprintf("prometheus-%v", clusterID)
}
