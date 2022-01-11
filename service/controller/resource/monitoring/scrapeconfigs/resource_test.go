package scrapeconfigs

import (
	"flag"
	"path"
	"path/filepath"
	"testing"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	k8sclientfake "github.com/giantswarm/k8sclient/v7/pkg/k8sclient/fake"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestAWSScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			TemplatePath: path,
			Provider:     "aws",
			Vault:        "vault1.some-installation.test",
			Installation: "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			cluster, err := key.ToCluster(v)
			if err != nil {
				return nil, microerror.Mask(err)
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
				var logger micrologger.Logger
				{
					c := micrologger.Config{}

					logger, err = micrologger.New(c)
					if err != nil {
						t.Fatal(err)
					}
				}

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

			return toData(v, config)
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
			Vault:        "vault1.some-installation.test",
			Installation: "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			cluster, err := key.ToCluster(v)
			if err != nil {
				return nil, microerror.Mask(err)
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
						"crt": []byte("crt"),
					}}
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
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(v1.SchemeBuilder),
				}
				k8sClient, err = k8sclientfake.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}
			config.K8sClient = k8sClient

			return toData(v, config)
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

func TestKVMScrapeconfigs(t *testing.T) {
	var testFunc unittest.TestFunc
	{
		path := path.Join(unittest.ProjectRoot(), templatePath)

		config := Config{
			AdditionalScrapeConfigs: additionalScrapeConfigs,
			TemplatePath:            path,
			Provider:                "kvm",
			Vault:                   "vault1.some-installation.test",
			Installation:            "test-installation",
		}
		testFunc = func(v interface{}) (interface{}, error) {
			cluster, err := key.ToCluster(v)
			if err != nil {
				return nil, microerror.Mask(err)
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
						"crt": []byte("crt"),
					}}
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
					Logger:        logger,
					SchemeBuilder: k8sclient.SchemeBuilder(v1.SchemeBuilder),
				}
				k8sClient, err = k8sclientfake.NewClients(c, secret)
				if err != nil {
					t.Fatal(err)
				}
			}
			config.K8sClient = k8sClient
			return toData(v, config)
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
