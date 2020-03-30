module github.com/giantswarm/prometheus-meta-operator

go 1.13

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/coreos/prometheus-operator v0.34.0
	github.com/giantswarm/apiextensions v0.0.0-20191213075442-71155aa0f5b7
	github.com/giantswarm/backoff v0.0.0-20190913091243-4dd491125192 // indirect
	github.com/giantswarm/k8sclient v0.0.0-20191213144452-f75fead2ae06
	github.com/giantswarm/microendpoint v0.0.0-20191121160659-e991deac2653
	github.com/giantswarm/microerror v0.0.0-20191011121515-e0ebc4ecf5a5
	github.com/giantswarm/microkit v0.0.0-20191023091504-429e22e73d3e
	github.com/giantswarm/micrologger v0.0.0-20191014091141-d866337f7393
	github.com/giantswarm/operatorkit v0.0.0-20191209140411-5d098618662e
	github.com/giantswarm/to v0.0.0-20191022113953-f2078541ec95 // indirect
	github.com/giantswarm/versionbundle v0.0.0-20191206123034-be95231628ae
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/juju/errgo v0.0.0-20140925100237-08cceb5d0b53 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9 // indirect
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // indirect
	sigs.k8s.io/cluster-api v0.2.8
	sigs.k8s.io/controller-runtime v0.4.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.16.4
