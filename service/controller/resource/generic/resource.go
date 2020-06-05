package generic

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Interface interface {
	Create(metav1.Object) (metav1.Object, error)
	Update(metav1.Object) (metav1.Object, error)
	Get(name string, options metav1.GetOptions) (metav1.Object, error)
	Delete(name string, options *metav1.DeleteOptions) error
}
type Config struct {
	Client Interface
	Logger micrologger.Logger
	Name   string
	ToCR   func(interface{}) (metav1.Object, error)
}

type Resource struct {
	client Interface
	logger micrologger.Logger
	name   string
	toCR   func(interface{}) (metav1.Object, error)
}

func New(config Config) (*Resource, error) {
	if config.Client == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Client must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}
	if config.ToCR == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ToCR must not be empty", config)
	}

	r := &Resource{
		client: config.Client,
		logger: config.Logger,
		name:   config.Name,
		toCR:   config.ToCR,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return r.name
}
