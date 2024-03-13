package alertmanagerconfig

import (
	"flag"
	"path/filepath"
	"testing"

	"golang.org/x/net/http/httpproxy"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestRenderingOfAlertmanagerNotificationTemplateWithLegacyMonitoring(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		config := Config{
			Installation:   "test-installation",
			GrafanaAddress: "https://grafana",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return renderNotificationTemplate(unittest.ProjectRoot(), config)
		}
	}

	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/notification-template/classic/" + flavor)
		if err != nil {
			t.Fatal(err)
		}

		c := unittest.Config{
			OutputDir:            outputDir,
			T:                    t,
			TestFunc:             testFunc,
			Flavor:               flavor,
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

func TestRenderingOfAlertmanagerNotificationTemplateWithMimirEnabled(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		config := Config{
			Installation:   "test-installation",
			GrafanaAddress: "https://grafana",
			MimirEnabled:   true,
			BaseDomain:     "prometheus.installation-prometheus.svc",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return renderNotificationTemplate(unittest.ProjectRoot(), config)
		}
	}

	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/notification-template/mimir-enabled/" + flavor)
		if err != nil {
			t.Fatal(err)
		}

		c := unittest.Config{
			OutputDir:            outputDir,
			T:                    t,
			TestFunc:             testFunc,
			Flavor:               flavor,
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

func TestRenderingOfAlertmanagerConfig(t *testing.T) {
	var testFunc unittest.TestFunc
	{

		proxyConfig := httpproxy.Config{}

		config := Config{
			GrafanaAddress: "https://grafana",
			Installation:   "test-installation",
			OpsgenieKey:    "opsgenie-key",
			Proxy:          proxyConfig.ProxyFunc(),
			Pipeline:       "testing",
			SlackApiURL:    "https://slack",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return renderAlertmanagerConfig(unittest.ProjectRoot(), config)
		}
	}

	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/alertmanager-config/" + flavor)
		if err != nil {
			t.Fatal(err)
		}

		c := unittest.Config{
			OutputDir:            outputDir,
			T:                    t,
			TestFunc:             testFunc,
			TestFuncReturnsBytes: true,
			Flavor:               flavor,
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
