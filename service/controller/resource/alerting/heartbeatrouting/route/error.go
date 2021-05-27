package route

import (
	"github.com/giantswarm/microerror"
)

var emptyRouteError = &microerror.Error{
	Kind: "emptyRouteError",
}

// IsEmptyRouteError asserts emptyRouteError.
func IsEmptyRouteError(err error) bool {
	return microerror.Cause(err) == emptyRouteError
}
