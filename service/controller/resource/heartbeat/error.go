package heartbeat

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
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

var apiNotFoundError = &microerror.Error{
	Kind: "apiNotFoundError",
}

func IsApiNotFoundError(err error) bool {
	if microerror.Cause(err) == apiNotFoundError {
		return true
	}

	apiErr, ok := err.(*client.ApiError)
	return ok && apiErr.StatusCode == http.StatusNotFound
}
