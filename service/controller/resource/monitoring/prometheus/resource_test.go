package prometheus

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	k8sclientfake "github.com/giantswarm/k8sclient/v7/pkg/k8sclient/fake"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/core/v1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	capiexp "sigs.k8s.io/cluster-api/exp/api/v1beta1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestPrometheus(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var k8sClient k8sclient.Interface
	{
		c := k8sclient.ClientsConfig{
			Logger: logger,
			SchemeBuilder: k8sclient.SchemeBuilder{
				v1.SchemeBuilder.AddToScheme,
				capi.AddToScheme,
				capiexp.AddToScheme,
			},
		}
		k8sClient, err = k8sclientfake.NewClients(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	config := Config{
		Address:            "http://prometheus/cluster",
		Customer:           "Giant Swarm",
		EvaluationInterval: "60s",
		Installation:       "test-installation",
		Pipeline:           "testing",
		Provider:           "provider",
		Region:             "onprem",
		ImageRepository:    "giantswarm/prometheus",
		LogLevel:           "debug",
		Registry:           "quay.io",
		RetentionDuration:  "2w",
		ScrapeInterval:     "60s",
		Version:            "v2.28.1",
		K8sClient:          k8sClient,
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			return toPrometheus(context.Background(), v, config)
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
