package prometheus

type Prometheus struct {
	Address     string
	BaseDomain  string
	Bastions    string
	Storage     PrometheusStorage
	Retention   PrometheusRetention
	RemoteWrite PrometheusRemoteWrite
}

type PrometheusStorage struct {
	CreatePVC string
	Size      string
}

type PrometheusRetention struct {
	Duration string
	Size     string
}

type PrometheusRemoteWrite struct {
	URL       string
	BasicAuth PrometheusRemoteWriteBasicAuth
}

type PrometheusRemoteWriteBasicAuth struct {
	Username string
	Password string
}
