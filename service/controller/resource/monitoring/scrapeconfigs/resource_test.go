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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

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
			k8sClient.CtrlClient().Scheme().AddKnownTypes(appsv1alpha1.SchemeGroupVersion, &appsv1alpha1.App{})

			config := Config{
				TemplatePath:       path,
				OrganizationReader: FakeReader{},
				Provider: cluster.Provider{
					Kind:   "aws",
					Flavor: "vintage",
				},
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Pipeline:     "test-pipeline",
				Region:       "eu-central-1",
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}
			resource, err := New(config)
			if err != nil {
				t.Fatal(err)
			}
			return resource.toData(context.Background(), v)
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

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)
		testFunc = func(v interface{}) (interface{}, error) {
			var apps = []*appsv1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz-observability-bundle",
						Namespace: "org-my-organization",
					},
					Status: appsv1alpha1.AppStatus{
						Version: "0.2.0",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{

						Name:      "kubernetes-observability-bundle",
						Namespace: "org-my-organization",
					},
					Status: appsv1alpha1.AppStatus{
						Version: "0.4.0",
					},
				},
			}

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

			k8sClient.CtrlClient().Scheme().AddKnownTypes(appsv1alpha1.SchemeGroupVersion, &appsv1alpha1.App{})
			for _, app := range apps {
				app := app // Create a new variable inside the loop and assign the value of app to it
				err = k8sClient.CtrlClient().Create(context.Background(), app)
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
				Pipeline:     "test-pipeline",
				Region:       "eu-central-1",
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}

			resource, err := New(config)
			if err != nil {
				t.Fatal(err)
			}
			return resource.toData(context.Background(), v)
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

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)
		testFunc = func(v interface{}) (interface{}, error) {

			var apps = []*appsv1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz-observability-bundle",
						Namespace: "org-my-organization",
					},
					Status: appsv1alpha1.AppStatus{
						Version: "0.2.0",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kubernetes-observability-bundle",
						Namespace: "org-my-organization",
					},
					Status: appsv1alpha1.AppStatus{
						Version: "0.4.0",
					},
				},
			}

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

			k8sClient.CtrlClient().Scheme().AddKnownTypes(appsv1alpha1.SchemeGroupVersion, &appsv1alpha1.App{})
			for _, app := range apps {
				err = k8sClient.CtrlClient().Create(context.Background(), app)
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
				Pipeline:     "test-pipeline",
				Region:       "eu-central-1",
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}

			resource, err := New(config)
			if err != nil {
				t.Fatal(err)
			}

			return resource.toData(context.Background(), v)
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

	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)
		testFunc = func(v interface{}) (interface{}, error) {

			var apps = []*appsv1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz-observability-bundle",
						Namespace: "org-my-organization",
					},
					Status: appsv1alpha1.AppStatus{
						Version: "0.2.0",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kubernetes-observability-bundle",
						Namespace: "org-my-organization",
					},
					Status: appsv1alpha1.AppStatus{
						Version: "0.4.0",
					},
				},
			}

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

			k8sClient.CtrlClient().Scheme().AddKnownTypes(appsv1alpha1.SchemeGroupVersion, &appsv1alpha1.App{})
			for _, app := range apps {
				err = k8sClient.CtrlClient().Create(context.Background(), app)
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
				Pipeline:     "test-pipeline",
				Region:       "eu-central-1",
				Customer:     "pmo",
				K8sClient:    k8sClient,
				Vault:        "vault1.some-installation.test",
				Installation: "test-installation",
				Logger:       logger,
			}

			resource, err := New(config)
			if err != nil {
				t.Fatal(err)
			}
			return resource.toData(context.Background(), v)
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
