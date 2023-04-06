package heartbeatwebhookconfig

import (
	"flag"
	"path/filepath"
	"testing"

	"golang.org/x/net/http/httpproxy"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestAlertmanager(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	proxyConfig := httpproxy.Config{}
	config := Config{
		Proxy:        proxyConfig.ProxyFunc(),
		Installation: "test-installation",
	}

	c := unittest.Config{
		OutputDir: outputDir,
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
