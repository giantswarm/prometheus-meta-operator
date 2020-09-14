package scrapeconfigs

import (
	"flag"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestScrapeconfigs(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot get current filename")
	}

	path := path.Join(path.Dir(filename), "../../../..", templatePath)

	config := Config{
		TemplatePath: path,
		Provider:     "aws",
		Vault:        "vault-address",
	}
	toScrapeConfigs := func(v interface{}) (metav1.Object, error) {
		return toSecret(v, config)
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc:  toScrapeConfigs,
		Update:    *update,
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
