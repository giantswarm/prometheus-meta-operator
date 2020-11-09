package prometheus

import (
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
		Provider:          "kvm",
		Region:            "onprem",
		StorageSize:       "50Gi",
		RetentionDuration: "2w",
		RetentionSize:     "45Gi",
		RemoteWriteURL:    "http://grafana/api/prom/push",
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			return toPrometheus(v, config)
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
