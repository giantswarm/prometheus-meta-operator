package heartbeat

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	if r.mimirEnabled {
		r.logger.Debugf(ctx, "mimir is enabled, deleting")
		return r.EnsureDeleted(ctx, obj)
	}

	desired, err := toHeartbeat(obj, r.installation, r.pipeline)
	if err != nil {
		return microerror.Mask(err)
	}

	// By default, we create the heartbeat unless is it already configured and does not need to be updated;
	var heartbeatHasToBeCreated bool = true

	r.logger.Debugf(ctx, "checking if heartbeat is already configured")
	getResult, err := r.heartbeatClient.Get(ctx, desired.Name)
	heartbeatAlreadyConfigured := !IsApiNotFoundError(err)

	if heartbeatAlreadyConfigured && err != nil {
		r.logger.Debugf(ctx, "heartbeat is configured but could not be fetched")
		return microerror.Mask(err)
	}

	if heartbeatAlreadyConfigured {
		// We need to delete and recreate it because the update is a PATCH (so existing alert tags are kept)
		// This causes issue when installations are switched from the testing pipeline to the stable pipeline as heartbeat are skipped.
		r.logger.Debugf(ctx, "heartbeat is configured")
		r.logger.Debugf(ctx, "checking if heartbeat needs to be reconfigured")

		var current heartbeat.Heartbeat = getResult.Heartbeat
		if hasChanged(current, *desired) {
			r.logger.Debugf(ctx, "heartbeat has changed and needs to be reconfigured")

			r.logger.Debugf(ctx, "deleting heartbeat")
			_, err := r.heartbeatClient.Delete(ctx, desired.Name)
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Debugf(ctx, "deleted heartbeat")
		} else {
			heartbeatHasToBeCreated = false
			r.logger.Debugf(ctx, "heartbeat is up to date")
		}
	}

	if heartbeatHasToBeCreated {
		r.logger.Debugf(ctx, "creating heartbeat")
		err := r.createAndPingHeartbeat(ctx, desired)
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "created heartbeat")
	}

	return nil
}

func (r *Resource) createAndPingHeartbeat(ctx context.Context, h *heartbeat.Heartbeat) error {
	req := &heartbeat.AddRequest{
		Name:          h.Name,
		Description:   h.Description,
		Interval:      h.Interval,
		IntervalUnit:  heartbeat.Unit(h.IntervalUnit),
		Enabled:       &h.Enabled,
		OwnerTeam:     h.OwnerTeam,
		AlertMessage:  h.AlertMessage,
		AlertTag:      h.AlertTags,
		AlertPriority: h.AlertPriority,
	}
	_, err := r.heartbeatClient.Add(ctx, req)
	if err != nil {
		return microerror.Mask(err)
	}

	// The initial ping to the heartbeat is there to move the heartbeat from inactive to active.
	_, err = r.heartbeatClient.Ping(ctx, h.Name)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
