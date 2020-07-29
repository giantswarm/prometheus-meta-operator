package prometheus

type Prometheus struct {
	BaseDomain string
	Storage    PrometheusStorage
}

type PrometheusStorage struct {
	CreatePVC string
	Size      string
}
