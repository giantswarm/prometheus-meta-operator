package pvresizing

const (
	Name = "prometheus"

	PrometheusStorageSizeSmall  PrometheusStorageSizeType = "small"
	PrometheusStorageSizeMedium PrometheusStorageSizeType = "medium"
	PrometheusStorageSizeLarge  PrometheusStorageSizeType = "large"
)

type PrometheusStorageSizeType string

// PrometheusDiskSize returns the desired disk size based on the
// value of annotation monitoring.giantswarm.io/prometheus-disk-size
func PrometheusDiskSize(v string) string {
	value := PrometheusStorageSizeType(v)
	switch value {
	case PrometheusStorageSizeSmall:
		return "30Gi"
	case PrometheusStorageSizeMedium:
		return "100Gi"
	case PrometheusStorageSizeLarge:
		return "200Gi"
	default:
		return "100Gi"
	}
}
