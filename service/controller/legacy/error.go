package legacy

import (
	"github.com/giantswarm/microerror"
)

var invalidProviderError = &microerror.Error{
	Kind: "invalidProviderError",
}

// IsInvalidProvider asserts invalidProviderError.
func IsInvalidProvider(err error) bool {
	return microerror.Cause(err) == invalidProviderError
}
