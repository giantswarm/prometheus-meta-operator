package project

var (
	description = "The prometheus-meta-operator does something."
	gitSHA      = "n/a"
	name        = "prometheus-meta-operator"
	source      = "https://github.com/giantswarm/prometheus-meta-operator"
	version     = "4.78.0"
)

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
