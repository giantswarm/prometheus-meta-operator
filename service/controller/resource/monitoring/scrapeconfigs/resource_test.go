package scrapeconfigs

import (
	"context"
	"flag"
	"path"
	"path/filepath"
	"testing"

	appsv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	fakek8sclient "github.com/giantswarm/k8sclient/v7/pkg/k8sclient/fake"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

var update = flag.Bool("update", false, "update the output file")

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

type FakeReader struct{}

func (r FakeReader) Read(ctx context.Context, cluster metav1.Object) (string, error) {
	return "my-organization", nil
}

func TestAWSScrapeconfigs(t *testing.T) {
	var err error
	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		var client client.Client
		{
			schemeBuilder := runtime.SchemeBuilder(k8sclient.SchemeBuilder{
				apiextensionsv1.AddToScheme,
				appsv1alpha1.AddToScheme,
			})

			err := schemeBuilder.AddToScheme(scheme.Scheme)
			if err != nil {
				t.Fatal(err)
			}

			client = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithRuntimeObjects().
				Build()
		}

		testFunc = func(v interface{}) (interface{}, error) {
			testCluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}
			var secret runtime.Object
			{
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-certificates",
						Namespace: key.Namespace(testCluster),
					},
					Data: map[string][]byte{
						"token": []byte("token"),
					},
				}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(corev1.SchemeBuilder),
				}
				k8sClient, err = fakek8sclient.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}

			config := Config{
				TemplatePath:       path,
				OrganizationReader: FakeReader{},
				Provider: cluster.Provider{
					Kind:   "aws",
					Flavor: "vintage",
				},
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/aws")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		Flavor:               "vintage",
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		TestFuncReturnsBytes: true,
		Update:               *update,
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
	var err error
	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		var client client.Client
		{
			schemeBuilder := runtime.SchemeBuilder(k8sclient.SchemeBuilder{
				apiextensionsv1.AddToScheme,
				appsv1alpha1.AddToScheme,
			})

			err := schemeBuilder.AddToScheme(scheme.Scheme)
			if err != nil {
				t.Fatal(err)
			}

			client = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithRuntimeObjects().
				Build()
		}

		testFunc = func(v interface{}) (interface{}, error) {
			testCluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}
			var secret runtime.Object
			{
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-certificates",
						Namespace: key.Namespace(testCluster),
					},
					Data: map[string][]byte{
						"crt": []byte("crt"),
						"key": []byte("key"),
					},
				}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(corev1.SchemeBuilder),
				}
				k8sClient, err = fakek8sclient.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}

			config := Config{
				TemplatePath:       path,
				OrganizationReader: FakeReader{},
				Provider: cluster.Provider{
					Kind:   "azure",
					Flavor: "vintage",
				},
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/azure")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		Flavor:               "vintage",
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		TestFuncReturnsBytes: true,
		Update:               *update,
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

func TestCAPZScrapeconfigs(t *testing.T) {
	var err error
	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var apps = []runtime.Object{
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "baz-observability-bundle",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.2.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "kubernetes-observability-bundle",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.4.0",
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
		testFunc = func(v interface{}) (interface{}, error) {
			testCluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}
			var secret runtime.Object
			{
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-certificates",
						Namespace: key.Namespace(testCluster),
					},
					Data: map[string][]byte{
						"crt": []byte("crt"),
						"key": []byte("key"),
					},
				}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(corev1.SchemeBuilder),
				}
				k8sClient, err = fakek8sclient.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}

			config := Config{
				AdditionalScrapeConfigs: additionalScrapeConfigs,
				TemplatePath:            path,
				OrganizationReader:      FakeReader{},
				Provider: cluster.Provider{
					Kind:   "capz",
					Flavor: "capi",
				},
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/capz")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		Flavor:               "capi",
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		TestFuncReturnsBytes: true,
		Update:               *update,
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
	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var apps = []runtime.Object{
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "baz-observability-bundle",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.2.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "kubernetes-observability-bundle",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.4.0",
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
		testFunc = func(v interface{}) (interface{}, error) {
			testCluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}
			var secret runtime.Object
			{
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-certificates",
						Namespace: key.Namespace(testCluster),
					},
					Data: map[string][]byte{
						"crt": []byte("crt"),
						"key": []byte("key"),
					},
				}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(corev1.SchemeBuilder),
				}
				k8sClient, err = fakek8sclient.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}

			config := Config{
				AdditionalScrapeConfigs: additionalScrapeConfigs,
				TemplatePath:            path,
				OrganizationReader:      FakeReader{},
				Provider: cluster.Provider{
					Kind:   "gcp",
					Flavor: "capi",
				},
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/gcp")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		Flavor:               "capi",
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		TestFuncReturnsBytes: true,
		Update:               *update,
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
	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var apps = []runtime.Object{
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "baz-observability-bundle",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.2.0",
			},
		},
		&appsv1alpha1.App{
			ObjectMeta: v1.ObjectMeta{
				Name:      "kubernetes-observability-bundle",
				Namespace: "org-my-organization",
			},
			Status: appsv1alpha1.AppStatus{
				Version: "0.4.0",
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
		testFunc = func(v interface{}) (interface{}, error) {
			testCluster, err := key.ToCluster(v)
			if err != nil {
				t.Fatal(err)
			}
			var secret runtime.Object
			{
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cluster-certificates",
						Namespace: key.Namespace(testCluster),
					},
					Data: map[string][]byte{
						"crt": []byte("crt"),
						"key": []byte("key"),
					},
				}
			}

			var k8sClient k8sclient.Interface
			{
				c := k8sclient.ClientsConfig{
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(corev1.SchemeBuilder),
				}
				k8sClient, err = fakek8sclient.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}

			config := Config{
				AdditionalScrapeConfigs: additionalScrapeConfigs,
				TemplatePath:            path,
				OrganizationReader:      FakeReader{},
				Provider: cluster.Provider{
					Kind:   "capa",
					Flavor: "capi",
				},
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}
			return toData(context.Background(), client, v, config)
		}
	}

	outputDir, err := filepath.Abs("./test/capa")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		Flavor:               "capi",
		OutputDir:            outputDir,
		T:                    t,
		TestFunc:             testFunc,
		TestFuncReturnsBytes: true,
		Update:               *update,
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
