package pvcresizing

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

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

func GetRetentionSize(storageSize resource.Quantity) string {
	// Set Retention.Size (TSDB limit) to a ratio of the volume storage size.
	storageSize.Set(int64(storageSize.AsApproximateFloat64() * VOLUME_STORAGE_LIMIT))
	return storageSize.String()
}
