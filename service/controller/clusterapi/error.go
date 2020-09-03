package clusterapi

import "github.com/giantswarm/microerror"

var unsupportedStorageVersionError = &microerror.Error{
	Kind: "unsupportedStorageVersionError",
}

// IsUnsupportedStorageVersion asserts unsupportedStorageVersionError.
func IsUnsupportedStorageVersion(err error) bool {
	return microerror.Cause(err) == unsupportedStorageVersionError
}
