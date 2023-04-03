package prometheus

type Prometheus struct {
	AdditionalScrapeConfigs string
	Address                 string
	BaseDomain              string
	Bastions                string
	EvaluationInterval      string
	LogLevel                string
	Mayu                    string
	Retention               PrometheusRetention
	ScrapeInterval          string
	Version                 string
}

type PrometheusRetention struct {
	Duration string
}

type PrometheusRemoteWriteBasicAuth struct {
	Username string
	Password string
}
