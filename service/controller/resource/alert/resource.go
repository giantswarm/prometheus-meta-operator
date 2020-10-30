package alert

import (
	"bytes"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/pkg/template"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name              = "alert"
	rulesFileLocation = "/opt/prometheus-meta-operator/files/templates/rules/**/*.yml"
)

type Config struct {
	Installation     string
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
}

type Resource struct {
	prometheusClient promclient.Interface
	logger           micrologger.Logger
	installation     string
}

func New(config Config) (*Resource, error) {
	if config.PrometheusClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	r := &Resource{
		logger:           config.Logger,
		prometheusClient: config.PrometheusClient,
		installation:     config.Installation,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return "alert"
}

type TemplateData struct {
	ClusterID    string
	Installation string
	ManagedBy    string
	Namespace    string
}

func getRules(obj interface{}, installation string) ([]*promv1.PrometheusRule, error) {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var data TemplateData = TemplateData{
		ClusterID:    key.ClusterID(cluster),
		Installation: installation,
		ManagedBy:    project.Name(),
		Namespace:    key.Namespace(cluster),
	}

	template, err := template.RenderTemplate(data, rulesFileLocation)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Split multi yaml files
	files := bytes.Split(template, []byte("---"))

	var rules []*promv1.PrometheusRule = make([]*promv1.PrometheusRule, 0)
	for _, file := range files {
		if len(bytes.TrimSpace(file)) > 0 {
			var rule promv1.PrometheusRule = promv1.PrometheusRule{}

			if err = yaml.Unmarshal(file, &rule); err != nil {
				return nil, microerror.Mask(err)
			}

			rules = append(rules, &rule)
		}
	}

	return rules, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.PrometheusRule)
	d := desired.(*promv1.PrometheusRule)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetLabels(c.GetLabels())
	d.SetAnnotations(c.GetAnnotations())
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetClusterName(c.GetClusterName())
	d.SetManagedFields(c.GetManagedFields())
}
