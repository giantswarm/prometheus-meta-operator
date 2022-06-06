package promremotewrite

import "github.com/giantswarm/microerror"

var errorFetchingPrometheus = &microerror.Error{
	Kind: "errorFetchingPrometheus",
}

var noSuchPrometheusForLabel = &microerror.Error{
	Kind: "noSuchPrometheusForLabel",
}

// IsNoSuchPrometheusForLabel asserts noSuchPrometheusForLabel.
func IsNoSuchPrometheusForLabel(err error) bool {
	return microerror.Cause(err) == noSuchPrometheusForLabel
}

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongType asserts wrongTypeError.
func IsWrongType(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
