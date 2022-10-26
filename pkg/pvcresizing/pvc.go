package pvcresizing

const (
	PrometheusStorageSizeSmall  PrometheusStorageSizeType = "small"
	PrometheusStorageSizeMedium PrometheusStorageSizeType = "medium"
	PrometheusStorageSizeLarge  PrometheusStorageSizeType = "large"
)

type PrometheusStorageSizeType string

// PrometheusVolumeSize returns the desired disk size based on the
// param annotationValue: small, medium Or large
func PrometheusVolumeSize(annotationValue string) string {
	value := PrometheusStorageSizeType(annotationValue)
	switch value {
	case PrometheusStorageSizeSmall:
		return "30Gi"
	case PrometheusStorageSizeLarge:
		return "200Gi"
	default:
		return "100Gi"
	}
}
