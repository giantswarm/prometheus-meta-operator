package alertmanagerconfig

import (
	"flag"
	"path"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestAlertmanagerconfigsecret(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			TemplatePath:     path,
			Provider:         "aws",
			Installation:     "test-installation",
			OpsgenieKey:      "opsgenie-key",
			GrafanaAddress:   "https://grafana",
			SlackApiURL:      "https://slack",
			SlackProjectName: "my-slack-project",
			Pipeline:         "testing",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(v, config)
		}
	}

	outputDir, err := filepath.Abs("./test")
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
