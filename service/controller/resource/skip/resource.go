package skip

import (
	"github.com/giantswarm/micrologger"
)

const (
	Name = "skip"
)

type Config struct {
	Logger       micrologger.Logger
	Installation string
}

type Resource struct {
	logger       micrologger.Logger
	installation string
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		logger:       config.Logger,
		installation: config.Installation,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
