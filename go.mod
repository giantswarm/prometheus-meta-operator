module github.com/giantswarm/prometheus-meta-operator

go 1.13

require (
	github.com/coreos/prometheus-operator v0.39.0
	github.com/giantswarm/apiextensions v0.4.6
	github.com/giantswarm/k8sclient/v3 v3.1.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/microkit v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit v1.0.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/spf13/viper v1.6.2
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/cluster-api v0.2.8
)

replace k8s.io/client-go => k8s.io/client-go v0.16.4

replace github.com/giantswarm/operatorkit => ./operatorkit
