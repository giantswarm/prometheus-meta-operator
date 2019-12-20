package servicemonitor

import (
	"fmt"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

func toServiceMonitors(obj interface{}) ([]*promv1.ServiceMonitor, error) {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return []*promv1.ServiceMonitor{
		apiServer(cluster),
	}, nil
}

func apiServer(cluster *v1alpha2.Cluster) *promv1.ServiceMonitor {
	return &promv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("kubernetes-apiserver-%s", cluster.GetName()),
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
		},
		Spec: promv1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"component": "apiserver",
				},
			},
			NamespaceSelector: promv1.NamespaceSelector{
				Any: true,
			},
			Endpoints: []promv1.Endpoint{
				promv1.Endpoint{
					Port:   "https",
					Scheme: "https",
					TLSConfig: &promv1.TLSConfig{
						CAFile:             fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.Secret()),
						CertFile:           fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.Secret()),
						KeyFile:            fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.Secret()),
						InsecureSkipVerify: true,
					},
				},
			},
		},
	}
}
