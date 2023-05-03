package remotewritesecret

import (
	"github.com/giantswarm/microerror"
)

var remoteWriteNotFound = &microerror.Error{
	Kind: "remoteWriteNotFound",
}

// IsRemoteWriteNotFound asserts remoteWriteNotFound.
func IsRemoteWriteNotFound(err error) bool {
	return microerror.Cause(err) == remoteWriteNotFound
}
