package alertmanager

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestAlertmanager(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		Address:     "http://alertmanager/cluster",
		StorageSize: "1Gi",
		Version:     "v0.22.1",
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			return toAlertmanager(context.TODO(), v, config)
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
