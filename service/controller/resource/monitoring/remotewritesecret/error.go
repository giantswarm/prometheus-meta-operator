package remotewritesecret

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

var remoteWriteNotFound = &microerror.Error{
	Kind: "remoteWriteNotFound",
}

// IsRemoteWriteNotFound asserts remoteWriteNotFound.
func IsRemoteWriteNotFound(err error) bool {
	return microerror.Cause(err) == remoteWriteNotFound
}
