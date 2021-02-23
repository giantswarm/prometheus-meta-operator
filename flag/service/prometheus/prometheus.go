package prometheus

type Prometheus struct {
	Address     string
	BaseDomain  string
	Bastions    string
	Mayu        string
	Storage     PrometheusStorage
	Retention   PrometheusRetention
	RemoteWrite PrometheusRemoteWrite
	Version     string
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
