package prometheus

type Prometheus struct {
	BaseDomain string
	Storage    PrometheusStorage
}

type PrometheusStorage struct {
	IsPersistent string
	Size         string
}
