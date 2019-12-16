package project

var (
	bundleVersion = "0.0.1"
	description   = "The template-operator does something."
	gitSHA        = "n/a"
	name          = "template-operator"
	source        = "https://github.com/giantswarm/template-operator"
	version       = "n/a"
)

func BundleVersion() string {
	return bundleVersion
}

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
