package heartbeat

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	desired, err := toHeartbeat(obj, r.installation)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if heartbeat exists")
	var current heartbeat.Heartbeat
	getResult, err := r.heartbeatClient.Get(ctx, desired.Name)
	if IsApiNotFoundError(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "heartbeat dot not exists")
		r.logger.LogCtx(ctx, "level", "debug", "message", "creating heartbeat")
		req := &heartbeat.AddRequest{
			Name:          desired.Name,
			Description:   desired.Description,
			Interval:      desired.Interval,
			IntervalUnit:  heartbeat.Unit(desired.IntervalUnit),
			Enabled:       &desired.Enabled,
			OwnerTeam:     desired.OwnerTeam,
			AlertMessage:  desired.AlertMessage,
			AlertTag:      desired.AlertTags,
			AlertPriority: desired.AlertPriority,
		}
		addResult, err := r.heartbeatClient.Add(ctx, req)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "created heartbeat")
		current = addResult.Heartbeat
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "heartbeat exists")
		current = getResult.Heartbeat
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if heartbeat needs update")
	if hasChanged(current, *desired) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "heartbeat needs update")
		r.logger.LogCtx(ctx, "level", "debug", "message", "updating heartbeat")
		req := &heartbeat.UpdateRequest{
			Name:          desired.Name,
			Description:   desired.Description,
			Interval:      desired.Interval,
			IntervalUnit:  heartbeat.Unit(desired.IntervalUnit),
			Enabled:       &desired.Enabled,
			OwnerTeam:     desired.OwnerTeam,
			AlertMessage:  desired.AlertMessage,
			AlertTag:      desired.AlertTags,
			AlertPriority: desired.AlertPriority,
		}
		_, err := r.heartbeatClient.Update(ctx, req)
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "updated heartbeat")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "heartbeat is up to date")
	}

	return nil
}
