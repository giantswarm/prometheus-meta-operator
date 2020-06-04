module github.com/giantswarm/prometheus-meta-operator

go 1.13

require (
	github.com/coreos/prometheus-operator v0.34.0
	github.com/giantswarm/apiextensions v0.4.6
	github.com/giantswarm/k8sclient/v3 v3.1.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/microkit v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit v1.0.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/spf13/viper v1.6.2
	k8s.io/api v0.16.6
	k8s.io/apimachinery v0.16.6
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // indirect
	sigs.k8s.io/cluster-api v0.2.8
)

replace k8s.io/client-go => k8s.io/client-go v0.16.4

replace github.com/giantswarm/operatorkit => ./operatorkit
