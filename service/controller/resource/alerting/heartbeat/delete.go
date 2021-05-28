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

	r.logger.Debugf(ctx, "checking if heartbeat exists")
	_, err = r.heartbeatClient.Get(ctx, desired.Name)
	if IsApiNotFoundError(err) {
		r.logger.Debugf(ctx, "heartbeat does not exist")
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		r.logger.Debugf(ctx, "triggering final heartbeat ping")
		// The final ping to the heartbeat cleans up any opened heartbeat alerts for the cluster being deleted.
		_, err = r.heartbeatClient.Ping(ctx, desired.Name)
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "triggered final heartbeat ping")

		r.logger.Debugf(ctx, "deleting heartbeat")
		_, err = r.heartbeatClient.Delete(ctx, desired.Name)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Debugf(ctx, "deleted heartbeat")
	}

	return nil
}
