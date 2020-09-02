module github.com/giantswarm/prometheus-meta-operator

go 1.13

require (
	github.com/coreos/prometheus-operator v0.41.1
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/apiextensions/v2 v2.1.0
	github.com/giantswarm/k8sclient/v4 v4.0.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/microkit v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit/v2 v2.0.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-cmp v0.5.1
	github.com/spf13/viper v1.6.2
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	k8s.io/api v0.18.5
	k8s.io/apiextensions-apiserver v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	sigs.k8s.io/cluster-api v0.3.7
	sigs.k8s.io/controller-runtime v0.6.1
)
