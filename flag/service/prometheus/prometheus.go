package prometheus

type Prometheus struct {
	AdditionalScrapeConfigs string
	Address                 string
	BaseDomain              string
	Bastions                string
	EvaluationInterval      string
	LogLevel                string
	Mayu                    string
	Storage                 PrometheusStorage
	Retention               PrometheusRetention
	ScrapeInterval          string
	Version                 string
}

type PrometheusStorage struct {
	Size string
}

type PrometheusRetention struct {
	Duration string
	Size     string
}

type PrometheusRemoteWriteBasicAuth struct {
	Username string
	Password string
}
