module github.com/giantswarm/prometheus-meta-operator

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/giantswarm/apiextensions/v3 v3.22.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/k8sclient/v5 v5.11.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit/v4 v4.3.1
	github.com/giantswarm/versionbundle v0.2.0
	github.com/go-kit/kit v0.10.0
	github.com/go-logr/logr v0.3.0 // indirect
	github.com/go-openapi/errors v0.20.0 // indirect
	github.com/go-openapi/runtime v0.19.26 // indirect
	github.com/go-openapi/validate v0.20.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.5
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f
	github.com/opsgenie/opsgenie-go-sdk-v2 v1.2.8
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator v0.43.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.21.0
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/autoscaler/vertical-pod-autoscaler v0.9.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20201113171705-d219536bb9fd // indirect
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	sigs.k8s.io/cluster-api v0.3.15
	sigs.k8s.io/controller-runtime v0.6.4
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	// cf. https://github.com/prometheus/prometheus/issues/6048#issuecomment-534549253 and https://github.com/ctron/enmasse/commit/4e30e62e83fcfcd87b4cdada3b30edbffbe09e85
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.0.0-20201126101154-26d89b4b0776
	k8s.io/client-go => k8s.io/client-go v0.19.4
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
)
