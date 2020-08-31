# Targets

- cadvisor
- calico
- cert-exporter
- chart-operator
- cluster-autoscaler
- docker
- etcd
- kube-proxy
- kube-state-metrics
- kubelet
- kubernetes-apiserver
- kubernetes-controller-manager
- kubernetes-scheduler
- net-exporter
- nginx-ingress-controller
- node-exporter

### AWS specific

- aws-node

## Known issues

* docker: fails to be scrapped, due to metrics port not exposed on giantswarm releases below 11.0.0, see [commit](https://github.com/giantswarm/k8scloudconfig/commit/6ecc07e665c3e854dfa8be102a8c6446d1d9dc3c#diff-be6122463e3fe598d118a80e09254d3d)

* metrics-server: not scrapping because according to official documentation it should not be used [as a source of monitoring solution metrics](https://github.com/kubernetes-sigs/metrics-server#kubernetes-metrics-server).

* nginx-ingress-controller: fails to be scrapped, due to network policy not allowing metrics port (10254). This is fix from [release v10.0.0](https://github.com/giantswarm/releases/tree/master/aws/archived/v10.1.0)
