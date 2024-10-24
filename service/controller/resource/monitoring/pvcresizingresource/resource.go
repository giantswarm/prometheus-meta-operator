package pvcresizingresource

import (
	"fmt"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/micrologger"
)

const (
	Name = "pvcresizingresource"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func resourceName(clusterID string) string {
	return fmt.Sprintf("prometheus-%v", clusterID)
}
