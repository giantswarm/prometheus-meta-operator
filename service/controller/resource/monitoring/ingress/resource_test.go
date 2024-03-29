package ingress

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestIngressDefault(t *testing.T) {
	testFunc := func(v interface{}) (interface{}, error) {
		return toIngress(v, Config{BaseDomain: "prometheus"})
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

func TestIngressRestrictedAccess(t *testing.T) {
	testFunc := func(v interface{}) (interface{}, error) {
		return toIngress(v, Config{BaseDomain: "prometheus", RestrictedAccessEnabled: true, WhitelistedSubnets: "21.10.178/24"})
	}
	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/restricted-access/" + flavor)
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
		return toIngress(v, Config{BaseDomain: "prometheus", ExternalDNS: true})
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

func TestIngressExternalDNSWithRestrictedAccess(t *testing.T) {
	testFunc := func(v interface{}) (interface{}, error) {
		return toIngress(v, Config{BaseDomain: "prometheus.3lkdj.test.gigantic.io", ExternalDNS: true, RestrictedAccessEnabled: true, WhitelistedSubnets: "21.10.178/24"})
	}
	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/externaldns-with-restricted-access/" + flavor)
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
