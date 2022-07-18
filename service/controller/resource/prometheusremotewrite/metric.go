package prometheusremotewrite

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	reconcileErrors = prom.NewCounterVec(prom.CounterOpts{
		Name: "pmo_remotewrite_reconcile_errors_total",
		Help: "Total number of reconciliation errors for remotewrite controller",
	}, []string{"controller"})
)

func init() {
	prom.MustRegister(reconcileErrors)
}
