package prometheus

type Prometheus struct {
	BaseDomain string
	Security   PrometheusSecurity
}

type PrometheusSecurity struct {
	LetsEncrypt  PrometheusSecurityLetsEncrypt
	Whitelisting PrometheusSecurityWhitelisting
}

type PrometheusSecurityLetsEncrypt struct {
	Enabled string
}

type PrometheusSecurityWhitelisting struct {
	Enabled   string
	SourceIPs string
}
