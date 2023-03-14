package pvcresizing

import "strconv"

const (
	PrometheusStorageSizeSmall  PrometheusStorageSizeType = "small"
	PrometheusStorageSizeMedium PrometheusStorageSizeType = "medium"
	PrometheusStorageSizeLarge  PrometheusStorageSizeType = "large"

	// We apply a ratio to the volume storage size to compute the RetentionSize property (RetentionSize = 90% volume storage size)
	VOLUME_STORAGE_LIMIT = 0.9
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

func GetRetentionSize(storageSize float64) string {
	// Set Retention.Size (TSDB limit) to a ratio of the volume storage size.
	return strconv.FormatInt(int64(storageSize*VOLUME_STORAGE_LIMIT), 10)
}
