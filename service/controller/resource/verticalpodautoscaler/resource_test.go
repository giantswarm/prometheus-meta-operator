package verticalpodautoscaler

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_clientsetfake "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/fake"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestVerticalPodAutoScaler(t *testing.T) {
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

	var node runtime.Object
	{
		node = &v1.Node{
			Status: v1.NodeStatus{
				Allocatable: v1.ResourceList{
					v1.ResourceMemory: resource.MustParse("10"),
				},
			},
		}
	}

	var k8sClient *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger:        logger,
			SchemeBuilder: k8sclient.SchemeBuilder(v1.SchemeBuilder),
		}
		k8sClient, err = k8sclient.NewFakeClients(c, node)
		if err != nil {
			t.Fatal(err)
		}
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			c := Config{
				Logger:    logger,
				K8sClient: k8sClient,
				VpaClient: vpa_clientsetfake.NewSimpleClientset(),
			}
			r, err := New(c)
			if err != nil {
				return nil, err
			}
			return r.getObject(context.TODO(), v)
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
