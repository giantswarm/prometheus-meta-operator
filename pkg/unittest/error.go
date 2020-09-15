package unittest

import (
	"github.com/giantswarm/microerror"
)

var executionError = &microerror.Error{
	Kind: "executionError",
}

// IsExecution asserts executionError.
func IsExecution(err error) bool {
	return microerror.Cause(err) == executionError
}
