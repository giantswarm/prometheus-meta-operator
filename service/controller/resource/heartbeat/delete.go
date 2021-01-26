package heartbeat

import (
	"context"

	"github.com/giantswarm/microerror"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	desired, err := toHeartbeat(obj, r.installation, r.pipeline)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if heartbeat exists")
	_, err = r.heartbeatClient.Get(ctx, desired.Name)
	if IsApiNotFoundError(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "heartbeat does not exist")
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "triggering final heartbeat ping")
		// The final ping to the heartbeat cleans up any opened heartbeat alerts for the cluster being deleted.
		_, err = r.heartbeatClient.Ping(ctx, desired.Name)
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "triggered final heartbeat ping")

		r.logger.LogCtx(ctx, "level", "debug", "message", "deleting heartbeat")
		_, err = r.heartbeatClient.Delete(ctx, desired.Name)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "deleted heartbeat")
	}

	return nil
}
