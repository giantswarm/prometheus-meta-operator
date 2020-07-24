package scrapeconfigs

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name         = "scrapeconfigs"
	templatePath = "/opt/prometheus-meta-operator/files/templates/additional-scrape-configs.template"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:     clientFunc,
		Logger:         config.Logger,
		Name:           Name,
		ToCR:           toSecret,
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toSecret(v interface{}) (metav1.Object, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeconfigs, err := parseTemplate(cluster.GetName())
	if err != nil {
		return nil, microerror.Mask(err)
	}

	scrapeConfigsSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.PrometheusAdditionalScrapeConfigsSecretName(),
			Namespace: key.Namespace(cluster),
		},
		Data: map[string][]byte{
			key.PrometheusAdditionalScrapeConfigsName(): []byte(scrapeconfigs),
		},
		Type: "Opaque",
	}

	return scrapeConfigsSecret, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}

func parseTemplate(clusterName string) (string, error) {
	template, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", microerror.Mask(err)
	}

	var scrapeconfigs string = string(template)
	var apiServerUrl string = fmt.Sprintf("https://master.%s", clusterName)

	scrapeconfigs = strings.ReplaceAll(scrapeconfigs, "<SECRET_NAME>", key.Secret())
	scrapeconfigs = strings.ReplaceAll(scrapeconfigs, "<API_SERVER_URL>", apiServerUrl)

	return scrapeconfigs, nil
}
