package prometheus

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestPrometheus(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		Address:           "http://prometheus/cluster",
		CreatePVC:         true,
		Customer:          "Giant Swarm",
		Installation:      "test-installation",
		Pipeline:          "testing",
		Provider:          "provider",
		Region:            "onprem",
		LogLevel:          "debug",
		Registry:          "quay.io",
		StorageSize:       "50Gi",
		RetentionDuration: "2w",
		RetentionSize:     "45Gi",
		Version:           "v2.28.1",
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			return toPrometheus(context.TODO(), v, config)
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
