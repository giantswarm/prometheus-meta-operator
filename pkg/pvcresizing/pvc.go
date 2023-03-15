package pvcresizing

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

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

func GetRetentionSize(storageSize resource.Quantity) string {
	// Set Retention.Size (TSDB limit) to a ratio of the volume storage size.
	storageSize.Set(int64(storageSize.AsApproximateFloat64() * key.PrometheusVolumeStorageLimitRatio))
	// Trick to fit with the expected format in RetentionSize property
	// See Prometheus spec: https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#prometheusspec
	// See ByteSize format: https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#bytesizestring-alias
	return fmt.Sprintf("%sB", storageSize.String())
}
