package skip

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	if cluster.GetName() == r.installation {
		r.logger.Debugf(ctx, "cancel reconciliation")
		reconciliationcanceledcontext.SetCanceled(ctx)
		return nil
	}

	return nil
}
