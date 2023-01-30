package service

import (
	"github.com/giantswarm/operatorkit/v8/pkg/flag/service/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/alertmanager"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/grafana"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/installation"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/opsgenie"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/provider"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/security"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/slack"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/vault"
)

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Installation installation.Installation
	Kubernetes   kubernetes.Kubernetes
	Alertmanager alertmanager.Alertmanager
	Prometheus   prometheus.Prometheus
	Grafana      grafana.Grafana
	Slack        slack.Slack
	Opsgenie     opsgenie.Opsgenie
	Provider     provider.Provider
	Security     security.Security
	Vault        vault.Vault
}
