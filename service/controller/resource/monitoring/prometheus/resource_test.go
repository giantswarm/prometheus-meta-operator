package prometheus

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	k8sclientfake "github.com/giantswarm/k8sclient/v7/pkg/k8sclient/fake"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
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
		RemoteWriteURL:    "http://grafana/api/prom/push",
		Version:           "v2.28.1",
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			cluster, err := key.ToCluster(v)
			if err != nil {
				return nil, err
			}

			// Create a fake secret to get the authentication mechanism for this given cluster
			var secret runtime.Object
			{
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      key.SecretAPICertificates(cluster),
						Namespace: key.Namespace(cluster),
					},
					Data: map[string][]byte{
						"token": []byte("token"),
					}}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(v1.SchemeBuilder),
				}
				k8sClient, err = k8sclientfake.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}
			config.K8sClient = k8sClient

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
