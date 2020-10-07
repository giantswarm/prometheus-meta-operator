module github.com/giantswarm/prometheus-meta-operator

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
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
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/prometheus/common v0.13.0
	github.com/prometheus/prometheus v2.21.0+incompatible
	github.com/spf13/viper v1.6.2
	golang.org/x/crypto v0.0.0-20200930160638-afb6bcd081ae // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver v0.18.5
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	sigs.k8s.io/cluster-api v0.3.7
	sigs.k8s.io/controller-runtime v0.6.1
)
