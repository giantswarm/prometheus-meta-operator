package project

var (
	description = "The prometheus-meta-operator does something."
	gitSHA      = "n/a"
	name        = "prometheus-meta-operator"
	source      = "https://github.com/giantswarm/prometheus-meta-operator"
<<<<<<< HEAD
	version     = "4.64.1-dev"
=======
	version     = "4.64.0"
>>>>>>> e8a26219 (Release v4.64.0 (#1489))
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
