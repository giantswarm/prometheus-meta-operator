package project

import (
	"github.com/giantswarm/versionbundle"
)

func NewVersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   "template-operator",
				Description: "TODO",
				Kind:        versionbundle.KindChanged,
			},
		},
		Components: []versionbundle.Component{},
		Name:       "template-operator",
		Version:    BundleVersion(),
	}
}
