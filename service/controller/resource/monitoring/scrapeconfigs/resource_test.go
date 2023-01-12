package scrapeconfigs

import (
	"context"
	"flag"
	"path"
	"path/filepath"
	"testing"

	appsv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

const additionalScrapeConfigs = `- job_name: test1
  static_configs:
  - targets:
    - 1.1.1.1:123
  relabel_configs:
  - source_labels: [__address__]
    target_label: __param_target
- job_name: test2
  static_configs:
  - targets:
    - 8.8.8.8:123
  relabel_configs:
  - source_labels: [__address__]
    target_label: __param_target`

func TestAWSScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			TemplatePath: path,
			Provider:     "aws",
			Customer:     "pmo",
			Vault:        "vault1.some-installation.test",
			Installation: "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(context.Background(), nil, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/aws")
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

func TestAzureScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			TemplatePath: path,
			Provider:     "azure",
			Customer:     "pmo",
			Vault:        "vault1.some-installation.test",
			Installation: "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(context.Background(), nil, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/azure")
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

func TestKVMScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			AdditionalScrapeConfigs: additionalScrapeConfigs,
			TemplatePath:            path,
			Provider:                "kvm",
			Customer:                "pmo",
			Vault:                   "vault1.some-installation.test",
			Installation:            "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(context.Background(), nil, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/kvm")
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

func TestOpenStackScrapeconfigs(t *testing.T) {
	var err error

	var apps = []runtime.Object{
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "alice-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.7.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.7.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "bar-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.7.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "baz-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.8.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "kubernetes-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "1.0.0",
			},
		},
	}

	var client client.Client
	{
		schemeBuilder := runtime.SchemeBuilder(k8sclient.SchemeBuilder{
			apiextensionsv1.AddToScheme,
			appsv1alpha1.AddToScheme,
		})

		err = schemeBuilder.AddToScheme(scheme.Scheme)
		if err != nil {
			t.Fatal(err)
		}

		client = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithRuntimeObjects(apps...).
			Build()
	}

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			AdditionalScrapeConfigs: additionalScrapeConfigs,
			TemplatePath:            path,
			Provider:                "openstack",
			Customer:                "pmo",
			Vault:                   "vault1.some-installation.test",
			Installation:            "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/openstack")
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

func TestGCPScrapeconfigs(t *testing.T) {
	var err error

	var apps = []runtime.Object{
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "alice-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.16.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.15.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "bar-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.15.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "baz-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.16.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "kubernetes-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "1.0.0",
			},
		},
	}

	var client client.Client
	{
		schemeBuilder := runtime.SchemeBuilder(k8sclient.SchemeBuilder{
			apiextensionsv1.AddToScheme,
			appsv1alpha1.AddToScheme,
		})

		err = schemeBuilder.AddToScheme(scheme.Scheme)
		if err != nil {
			t.Fatal(err)
		}

		client = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithRuntimeObjects(apps...).
			Build()
	}

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			AdditionalScrapeConfigs: additionalScrapeConfigs,
			TemplatePath:            path,
			Provider:                "gcp",
			Customer:                "pmo",
			Vault:                   "vault1.some-installation.test",
			Installation:            "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/gcp")
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

func TestCAPAScrapeconfigs(t *testing.T) {
	var err error

	var apps = []runtime.Object{
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "alice-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.11.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.9.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "bar-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.10.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "baz-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.12.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "kubernetes-default-apps",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "1.0.0",
			},
		},
	}

	var client client.Client
	{
		schemeBuilder := runtime.SchemeBuilder(k8sclient.SchemeBuilder{
			apiextensionsv1.AddToScheme,
			appsv1alpha1.AddToScheme,
		})

		err = schemeBuilder.AddToScheme(scheme.Scheme)
		if err != nil {
			t.Fatal(err)
		}

		client = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithRuntimeObjects(apps...).
			Build()
	}

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			AdditionalScrapeConfigs: additionalScrapeConfigs,
			TemplatePath:            path,
			Provider:                "capa",
			Customer:                "pmo",
			Vault:                   "vault1.some-installation.test",
			Installation:            "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/capa")
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
