package prometheus

type Prometheus struct {
	AdditionalScrapeConfigs string
	Address                 string
	BaseDomain              string
	Bastions                string
	LogLevel                string
	Mayu                    string
	Storage                 PrometheusStorage
	Retention               PrometheusRetention
	Version                 string
}

type PrometheusStorage struct {
	CreatePVC string
	Size      string
}

type PrometheusRetention struct {
	Duration string
	Size     string
}

type PrometheusRemoteWriteBasicAuth struct {
	Username string
	Password string
}
