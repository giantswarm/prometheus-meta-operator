package prometheus

type Prometheus struct {
	Address    string
	BaseDomain string
	Storage    PrometheusStorage
}

type PrometheusStorage struct {
	CreatePVC string
	Size      string
}
