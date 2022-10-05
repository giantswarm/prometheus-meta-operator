package pvcresizingresource

import (
	"context"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	// No need to delete anything as PV gets deleted by prometheus-operator

	return nil
}
