package configuration

import (
	"github.com/giantswarm/microerror"
)

var secretNotFound = &microerror.Error{
	Kind: "secretNotFound",
}

// IsSecretNotFound asserts secretNotFound.
func IsSecretNotFound(err error) bool {
	return microerror.Cause(err) == secretNotFound
}

var remoteWriteNotFound = &microerror.Error{
	Kind: "remoteWriteNotFound",
}

// IsRemoteWriteNotFound asserts remoteWriteNotFound.
func IsRemoteWriteNotFound(err error) bool {
	return microerror.Cause(err) == remoteWriteNotFound
}
