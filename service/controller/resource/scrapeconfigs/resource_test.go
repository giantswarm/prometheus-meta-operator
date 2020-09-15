package scrapeconfigs

import (
	"flag"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestAWSScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatal("cannot get current filename")
		}

		path := path.Join(path.Dir(filename), "../../../..", templatePath)

		config := Config{
			TemplatePath: path,
			Provider:     "aws",
			Vault:        "vault1.some-installation.test",
		}
		testFunc = func(v interface{}) (metav1.Object, error) {
			return toSecret(v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/aws")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc:  testFunc,
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

func TestAzureScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatal("cannot get current filename")
		}

		path := path.Join(path.Dir(filename), "../../../..", templatePath)

		config := Config{
			TemplatePath: path,
			Provider:     "azure",
			Vault:        "vault1.some-installation.test",
		}
		testFunc = func(v interface{}) (metav1.Object, error) {
			return toSecret(v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/azure")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc:  testFunc,
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

func TestKVMScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			t.Fatal("cannot get current filename")
		}

		path := path.Join(path.Dir(filename), "../../../..", templatePath)

		config := Config{
			TemplatePath: path,
			Provider:     "kvm",
			Vault:        "vault1.some-installation.test",
		}
		testFunc = func(v interface{}) (metav1.Object, error) {
			return toSecret(v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/kvm")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc:  testFunc,
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
