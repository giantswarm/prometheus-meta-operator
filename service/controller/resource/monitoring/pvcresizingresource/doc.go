package pvcresizingresource

// This package is used to Watch Cluster CR and check if the value of
// the annotation `monitoring.giantswarm.io/prometheus-volume-size` has changed
// and triggers a PVC resize
// following this procedure https://github.com/prometheus-operator/prometheus-operator/issues/4079#issuecomment-1211989005
