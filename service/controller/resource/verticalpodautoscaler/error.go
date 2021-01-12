package verticalpodautoscaler

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var quantityConvertionError = &microerror.Error{
	Kind: "quantityConvertionError",
}

// IsQuantityConvertion asserts quantityConvertionError.
func IsQuantityConvertion(err error) bool {
	return microerror.Cause(err) == quantityConvertionError
}

var nodeMemoryNotFoundError = &microerror.Error{
	Kind: "nodeMemoryNotFoundError",
}

// IsNodeMemoryNotFound asserts nodeMemoryNotFoundError.
func IsNodeMemoryNotFound(err error) bool {
	return microerror.Cause(err) == nodeMemoryNotFoundError
}
