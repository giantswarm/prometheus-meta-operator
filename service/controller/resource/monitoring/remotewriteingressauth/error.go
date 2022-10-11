package remotewriteingressauth

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
