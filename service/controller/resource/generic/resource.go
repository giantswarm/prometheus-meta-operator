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

// Config contains dependencies for Resource struct.
type Config struct {
	// ClientFunc is a function that takes a namespace and returns a Kubernetes
	// client for resources of the specific kind in that namespace.
	ClientFunc func(string) Interface

	// Logger to use for logging.
	Logger micrologger.Logger

	// Name of the resource handler.
	// Note this is not the same thing as name of the resource.
	Name string

	// GetObjectMeta is a function that takes a resource object, casts it to
	// appropriate type and returns the metadata of that object, i.e. its
	// metav1.ObjectMeta part (name, namespace, labels, annotations, etc.).
	GetObjectMeta func(context.Context, interface{}) (metav1.ObjectMeta, error)

	// GetDesiredObject is a function that takes a resource object and returns the
	// object populated with the desired state.
	GetDesiredObject func(context.Context, interface{}) (metav1.Object, error)

	// HasChangedFunc is a function that takes two copies of an object - first
	// with existing state in the cluster and second with the desired state for
	// the given resource and returns true if there is a difference between
	// them and therefore state needs to be reconciled to match the desired.
	HasChangedFunc func(metav1.Object, metav1.Object) bool

	DeleteIfMimirEnabled bool
}

type Resource struct {
	clientFunc           func(string) Interface
	logger               micrologger.Logger
	name                 string
	getObjectMeta        func(context.Context, interface{}) (metav1.ObjectMeta, error)
	getDesiredObject     func(context.Context, interface{}) (metav1.Object, error)
	hasChangedFunc       func(metav1.Object, metav1.Object) bool
	deleteIfMimirEnabled bool
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
	if config.GetObjectMeta == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.GetObjectMeta must not be empty", config)
	}
	if config.GetDesiredObject == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ToCR must not be empty", config)
	}
	if config.HasChangedFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HasChangedFunc must not be empty", config)
	}

	r := &Resource{
		clientFunc:           config.ClientFunc,
		logger:               config.Logger,
		name:                 config.Name,
		getObjectMeta:        config.GetObjectMeta,
		getDesiredObject:     config.GetDesiredObject,
		hasChangedFunc:       config.HasChangedFunc,
		deleteIfMimirEnabled: config.DeleteIfMimirEnabled,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return r.name
}
