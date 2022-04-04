package etcdcertificates

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

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongType asserts wrongTypeError.
func IsWrongType(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}

var keyMissingError = &microerror.Error{
	Kind: "keyMissingError",
}

// IsKeyMissing asserts keyMissingError.
func IsKeyMissing(err error) bool {
	return microerror.Cause(err) == keyMissingError
}
