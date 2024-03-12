package noop

import (
	"context"

	"github.com/giantswarm/micrologger"
)

const (
	Name = "noop"
)

type Config struct {
	Logger micrologger.Logger
}

type Resource struct {
	logger micrologger.Logger
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		logger: config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "noop create")
	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "noop delete")
	return nil
}
