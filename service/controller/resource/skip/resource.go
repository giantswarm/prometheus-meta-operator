package skip

import (
	"github.com/giantswarm/micrologger"
)

const (
	Name = "skip"
)

type Config struct {
	Logger micrologger.Logger

	// Name will be checked against the current object name.
	// If those match the entire reconciliation loop will be canceled.
	Name string
}

type Resource struct {
	logger micrologger.Logger
	name   string
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		logger: config.Logger,
		name:   config.Name,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
