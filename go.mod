module github.com/giantswarm/prometheus-meta-operator

go 1.13

require (
	github.com/ant31/crd-validation v0.0.0-20180702145049-30f8a35d0ac2 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/brancz/gojsontoyaml v0.0.0-20190425155809-e8bd32d46b3d // indirect
	github.com/coreos/prometheus-operator v0.36.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/giantswarm/apiextensions v0.4.6
	github.com/giantswarm/k8sclient/v3 v3.1.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/microkit v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit v1.0.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/improbable-eng/thanos v0.3.2 // indirect
	github.com/kylelemons/godebug v0.0.0-20170820004349-d65d576e9348 // indirect
	github.com/mitchellh/hashstructure v0.0.0-20170609045927-2bca23e0e452 // indirect
	github.com/openshift/prom-label-proxy v0.1.1-0.20191016113035-b8153a7f39f1 // indirect
	github.com/prometheus/tsdb v0.8.0 // indirect
	github.com/spf13/viper v1.6.2
	k8s.io/api v0.16.6
	k8s.io/apimachinery v0.16.6
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // indirect
	sigs.k8s.io/cluster-api v0.2.8
)

replace k8s.io/client-go => k8s.io/client-go v0.16.4

replace github.com/giantswarm/operatorkit => ./operatorkit
