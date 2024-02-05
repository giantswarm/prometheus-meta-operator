package prometheus

type Prometheus struct {
	AdditionalScrapeConfigs string
	Address                 string
	BaseDomain              string
	Bastions                string
	EvaluationInterval      string
	ImageRepository         string
	LogLevel                string
	ScrapeInterval          string
	Version                 string
}
