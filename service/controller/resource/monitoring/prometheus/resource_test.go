package prometheus

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient/fake"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
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

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			cluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}
			var secret runtime.Object
			{
				secret = &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-certificates",
						Namespace: key.Namespace(cluster),
					},
					Data: map[string][]byte{
						"token": []byte("my-token"),
					},
				}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(v1.SchemeBuilder),
				}
				k8sClient, err = fake.NewClients(c, secret)
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
				K8sClient:          k8sClient,
				Provider:           "provider",
				Region:             "onprem",
				ImageRepository:    "giantswarm/prometheus",
				LogLevel:           "debug",
				Registry:           "quay.io",
				RetentionDuration:  "2w",
				ScrapeInterval:     "60s",
				Version:            "v2.28.1",
			}

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
