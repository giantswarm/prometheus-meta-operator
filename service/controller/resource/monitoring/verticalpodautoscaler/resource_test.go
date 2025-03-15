package verticalpodautoscaler

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	appsv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	k8sclientfake "github.com/giantswarm/k8sclient/v8/pkg/k8sclient/fake"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	vpa_clientsetfake "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/fake"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the output file")

func TestVerticalPodAutoScaler(t *testing.T) {
	var err error
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
					v1.ResourceCPU:    resource.MustParse("8"),
					v1.ResourceMemory: resource.MustParse("16Gi"),
				},
			},
		}
	}

	var k8sClient k8sclient.Interface
	{
		schemeBuilder := runtime.SchemeBuilder(k8sclient.SchemeBuilder{
			apiextensionsv1.AddToScheme,
			appsv1alpha1.AddToScheme,
		})

		err = schemeBuilder.AddToScheme(scheme.Scheme)
		if err != nil {
			t.Fatal(err)
		}
		c := k8sclient.ClientsConfig{
			Logger:        logger,
			SchemeBuilder: k8sclient.SchemeBuilder(schemeBuilder),
		}
		k8sClient, err = k8sclientfake.NewClients(c, node)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, flavor := range unittest.ProviderFlavors {
		outputDir, err := filepath.Abs("./test/" + flavor)
		if err != nil {
			t.Fatal(err)
		}

		testFunc := func(v interface{}) (interface{}, error) {
			c := Config{
				Logger:    logger,
				K8sClient: k8sClient,
				VpaClient: vpa_clientsetfake.NewSimpleClientset(),
				Provider: cluster.Provider{
					Kind:   "aws",
					Flavor: flavor,
				},
				Installation: "test-installation",
			}
			r, err := New(c)
			if err != nil {
				return nil, err
			}
			return r.getObject(context.Background(), v)
		}

		c := unittest.Config{
			Flavor:    flavor,
			OutputDir: outputDir,
			T:         t,
			TestFunc:  testFunc,
			Update:    *update,
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
}
