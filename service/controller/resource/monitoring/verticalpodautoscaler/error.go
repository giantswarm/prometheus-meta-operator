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

var quantityConversionError = &microerror.Error{
	Kind: "quantityConversionError",
}

// IsQuantityConversion asserts quantityConversionError.
func IsQuantityConversion(err error) bool {
	return microerror.Cause(err) == quantityConversionError
}

var nodeCpuNotFoundError = &microerror.Error{
	Kind: "nodeCpuNotFoundError",
}

// IsNodeCpuNotFound asserts nodeCpuNotFoundError.
func IsNodeCpuNotFound(err error) bool {
	return microerror.Cause(err) == nodeCpuNotFoundError
}

var nodeMemoryNotFoundError = &microerror.Error{
	Kind: "nodeMemoryNotFoundError",
}

// IsNodeMemoryNotFound asserts nodeMemoryNotFoundError.
func IsNodeMemoryNotFound(err error) bool {
	return microerror.Cause(err) == nodeMemoryNotFoundError
}
