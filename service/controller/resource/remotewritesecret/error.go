package remotewritesecret

import "github.com/giantswarm/microerror"

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

var errorFetchingPrometheus = &microerror.Error{
	Kind: "errorFetchingPrometheus",
}
