module github.com/giantswarm/prometheus-meta-operator

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/giantswarm/apiextensions/v3 v3.38.0
	github.com/giantswarm/k8sclient/v5 v5.12.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.4.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.6.0
	github.com/giantswarm/operatorkit/v4 v4.3.1
	github.com/giantswarm/versionbundle v0.2.0
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.0
	github.com/google/go-cmp v0.5.6
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f
	github.com/opsgenie/opsgenie-go-sdk-v2 v1.2.8
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.50.0
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.50.0
	github.com/prometheus/alertmanager v0.23.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.30.0
	github.com/prometheus/common/sigv4 v0.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.10.0
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver v0.19.4
	k8s.io/apimachinery v0.22.0
	k8s.io/autoscaler/vertical-pod-autoscaler v0.9.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	sigs.k8s.io/cluster-api v0.3.19
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	k8s.io/client-go => k8s.io/client-go v0.19.4
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
)
