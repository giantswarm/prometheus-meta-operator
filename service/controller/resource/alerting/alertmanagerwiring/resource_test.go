package alertmanagerwiring

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestAlertmanagerconfig(t *testing.T) {
	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/" + flavor)

		if err != nil {
			t.Fatal(err)
		}

		c := unittest.Config{
			Flavor:    flavor,
			OutputDir: outputDir,
			T:         t,
			TestFunc: func(v interface{}) (interface{}, error) {
				return toData(), nil
			},
			TestFuncReturnsBytes: true,
			Update:               *update,
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
