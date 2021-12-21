package project

import (
	"github.com/giantswarm/versionbundle"
)

func NewVersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Components: []versionbundle.Component{},
		Name:       "prometheus-meta-operator",
		Version:    Version(),
	}
}
