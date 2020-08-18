package generic

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Interface interface {
	Create(ctx context.Context, object metav1.Object, options metav1.CreateOptions) (metav1.Object, error)
	Update(ctx context.Context, object metav1.Object, options metav1.UpdateOptions) (metav1.Object, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (metav1.Object, error)
	Delete(ctx context.Context, name string, options *metav1.DeleteOptions) error
}

type Config struct {
	ClientFunc     func(string) Interface
	Logger         micrologger.Logger
	Name           string
	ToCR           func(interface{}) (metav1.Object, error)
	HasChangedFunc func(metav1.Object, metav1.Object) bool
}

type Resource struct {
	clientFunc     func(string) Interface
	logger         micrologger.Logger
	name           string
	toCR           func(interface{}) (metav1.Object, error)
	hasChangedFunc func(metav1.Object, metav1.Object) bool
}

func New(config Config) (*Resource, error) {
	if config.ClientFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ClientFunc must not be empty", config)
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
	if config.HasChangedFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HasChangedFunc must not be empty", config)
	}

	r := &Resource{
		clientFunc:     config.ClientFunc,
		logger:         config.Logger,
		name:           config.Name,
		toCR:           config.ToCR,
		hasChangedFunc: config.HasChangedFunc,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return r.name
}
