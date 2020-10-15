package prometheus

type Prometheus struct {
	Address    string
	BaseDomain string
	Bastions   string
	Storage    PrometheusStorage
	Retention  PrometheusRetention
}

type PrometheusStorage struct {
	CreatePVC string
	Size      string
}

type PrometheusRetention struct {
	Duration string
	Size     string
}
