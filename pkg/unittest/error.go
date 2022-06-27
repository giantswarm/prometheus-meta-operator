package unittest

import (
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
)

var executionError = &microerror.Error{
	Kind: "executionError",
}

// IsExecution asserts executionError.
func IsExecution(err error) bool {
	return microerror.Cause(err) == executionError
}

// EquateErrors returns true if the supplied errors are of the same type and
// produce identical strings. This mirrors the error comparison behaviour of
// https://github.com/go-test/deep,
//
// This differs from cmpopts.EquateErrors, which does not test for error strings
// and instead returns whether one error 'is' (in the errors.Is sense) the
// other.
func EquateErrors() cmp.Option {
	return cmp.Comparer(func(a, b error) bool {
		if a == nil || b == nil {
			return a == nil && b == nil
		}

		av := reflect.ValueOf(a)
		bv := reflect.ValueOf(b)
		if av.Type() != bv.Type() {
			return false
		}

		return a.Error() == b.Error()
	})
}
