package heartbeatwebhookconfig

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/micrologger"
	"golang.org/x/net/http/httpproxy"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestAlertmanagerConfig(t *testing.T) {
	var err error
	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	proxyConfig := httpproxy.Config{}
	config := Config{
		Proxy:        proxyConfig.ProxyFunc(),
		Logger:       logger,
		Installation: "test-installation",
	}

	resource, err := New(config)
	if err != nil {
		t.Fatal(err)
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
				return resource.toAlertmanagerConfig(v)
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
