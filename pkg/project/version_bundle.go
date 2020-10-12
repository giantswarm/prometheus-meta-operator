package project

import (
	"github.com/giantswarm/versionbundle"
)

func NewVersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   "prometheus-meta-operator",
				Description: "The prometheus-meta-operator manages Kubernetes clusters monitoring.",
				Kind:        versionbundle.KindChanged,
			},
		},
		Components: []versionbundle.Component{},
		Name:       "prometheus-meta-operator",
		Version:    BundleVersion(),
	}
}
