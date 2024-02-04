package key

import (
	"github.com/giantswarm/microerror"
)

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongType asserts wrongTypeError.
func IsWrongType(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}

var infrastructureRefNotFoundError = &microerror.Error{
	Kind: "infrastructureRefNotFoundError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInfrastructureRefNotFoundError(err error) bool {
	return microerror.Cause(err) == infrastructureRefNotFoundError
}
