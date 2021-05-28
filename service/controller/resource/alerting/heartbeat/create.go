package heartbeat

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	desired, err := toHeartbeat(obj, r.installation, r.pipeline)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "checking if heartbeat exists")
	var current heartbeat.Heartbeat
	getResult, err := r.heartbeatClient.Get(ctx, desired.Name)

	if IsApiNotFoundError(err) {
		r.logger.Debugf(ctx, "heartbeat does not exist")
		r.logger.Debugf(ctx, "creating heartbeat")
		err := r.createHeartbeat(ctx, desired)
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "created heartbeat")

		// The initial ping to the heartbeat is there to move the heartbeat from inactive to active.
		_, err = r.heartbeatClient.Ping(ctx, desired.Name)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "heartbeat exists")
	current = getResult.Heartbeat

	// We get the ID back from opsgenie so we update it in the heartbeat
	if desired.OwnerTeam.Name == current.OwnerTeam.Name {
		desired.OwnerTeam = current.OwnerTeam
	}

	r.logger.Debugf(ctx, "checking if heartbeat needs update")

	if hasChanged(current, *desired) {
		r.logger.Debugf(ctx, "heartbeat needs update")
		r.logger.Debugf(ctx, "updating heartbeat")
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
		r.logger.Debugf(ctx, "updated heartbeat")
	} else {
		r.logger.Debugf(ctx, "heartbeat is up to date")
	}

	return nil
}

func (r *Resource) createHeartbeat(ctx context.Context, h *heartbeat.Heartbeat) error {
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
	return nil
}
