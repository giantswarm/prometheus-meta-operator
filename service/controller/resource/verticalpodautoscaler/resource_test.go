package verticalpodautoscaler

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
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
	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			c := Config{
				Logger:    logger,
				K8sClient: k8sclient.NewFakeClients(),
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
