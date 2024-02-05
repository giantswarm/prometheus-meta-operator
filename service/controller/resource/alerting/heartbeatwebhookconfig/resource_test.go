package heartbeatwebhookconfig

import (
	"flag"
	"path/filepath"
	"testing"

	"golang.org/x/net/http/httpproxy"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestAlertmanager(t *testing.T) {
	proxyConfig := httpproxy.Config{}
	config := Config{
		Proxy:        proxyConfig.ProxyFunc(),
		Installation: "test-installation",
	}

	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/" + flavor)
		if err != nil {
			t.Fatal(err)
		}
		c := unittest.Config{
			OutputDir: outputDir,
			Flavor:    flavor,
			T:         t,
			TestFunc: func(v interface{}) (interface{}, error) {
				return toAlertmanagerConfig(v, config)
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
}
