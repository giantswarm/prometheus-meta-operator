package prometheusremotewrite

import "github.com/giantswarm/microerror"

var errorFetchingPrometheus = &microerror.Error{
	Kind: "errorFetchingPrometheus",
}

var errorRetrievingSecret = &microerror.Error{
	Kind: "errorRetrievingSecret",
}

var errorCreatingSecret = &microerror.Error{
	Kind: "errorCreatingSecret",
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
