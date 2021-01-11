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

var cannotConvertQuantityToInt64 = &microerror.Error{
	Kind: "cannotConvertQuantityToInt64",
}

// IsCannotConvertQuantityToInt64 asserts cannotConvertQuantityToInt64.
func IsCannotConvertQuantityToInt64(err error) bool {
	return microerror.Cause(err) == cannotConvertQuantityToInt64
}
