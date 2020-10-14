package heartbeat

import (
	"context"

	"github.com/giantswarm/microerror"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	desired, err := toHeartbeat(obj, r.installation)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if heartbeat exists")
	_, err = r.heartbeatClient.Get(ctx, desired.Name)
	if IsApiNotFoundError(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "heartbeat does not exists")
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "deleting heartbeat")
		_, err = r.heartbeatClient.Delete(ctx, desired.Name)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "deleted heartbeat")
	}

	return nil
}
