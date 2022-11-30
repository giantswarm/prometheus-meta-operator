package alertmanagerconfig

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestRenderingOfAlertmanagerNotificationTemplate(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		config := Config{
			GrafanaAddress: "https://grafana",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return renderNotificationTemplate(unittest.ProjectRoot(), config)
		}
	}

	outputDir, err := filepath.Abs("./test/notification-template")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		Update:               *update,
		TestFuncReturnsBytes: true,
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

func TestRenderingOfAlertmanagerConfig(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		config := Config{
			GrafanaAddress:   "https://grafana",
			Installation:     "test-installation",
			OpsgenieKey:      "opsgenie-key",
			Pipeline:         "testing",
			Provider:         "aws",
			SlackApiURL:      "https://slack",
			SlackProjectName: "my-slack-project",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return renderAlertmanagerConfig(unittest.ProjectRoot(), config)
		}
	}

	outputDir, err := filepath.Abs("./test/alertmanager-config")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		Update:               *update,
		TestFuncReturnsBytes: true,
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
