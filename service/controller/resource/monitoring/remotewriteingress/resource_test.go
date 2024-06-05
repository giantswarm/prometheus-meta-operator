package remotewriteingress

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestIngressDefault(t *testing.T) {
	testFunc := func(v interface{}) (interface{}, error) {
		resource, err := New(Config{
			BaseDomain: "prometheus",
		})
		if err != nil {
			t.Fatal(err)
		}
		return resource.toIngress(v)
	}
	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/default/" + flavor)
		if err != nil {
			t.Fatal(err)
		}

		c := unittest.Config{
			Flavor:    flavor,
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
}

func TestIngressExternalDNS(t *testing.T) {
	testFunc := func(v interface{}) (interface{}, error) {
		resource, err := New(Config{
			BaseDomain:  "prometheus",
			ExternalDNS: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		return resource.toIngress(v)
	}
	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/externaldns/" + flavor)
		if err != nil {
			t.Fatal(err)
		}

		c := unittest.Config{
			Flavor:    flavor,
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
}
