package promxy

import (
	"flag"
	"net/url"
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestServerGroup(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	url, _ := url.Parse("https://kubernetes.default:443")
	c := unittest.Config{
		OutputDir:  outputDir,
		T:          t,
		Marshaller: yaml.Marshal,
		TestFunc: func(v interface{}) (interface{}, error) {
			cluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}

			return toServerGroup(cluster, url, "test-installation", "kvm")
		},
		Update: *update,
	}
	runner, err := unittest.NewRunner(c)
	if err != nil {
		t.Fatal(err)
	}

	err = runner.Run()
	if err != nil {
		t.Fatal(err)
	}
}
