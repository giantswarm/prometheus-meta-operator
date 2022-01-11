package authentication

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func GetAPIAuthenticationMechanism(ctx context.Context, k8sclient kubernetes.Interface, cluster metav1.Object, clusterType string) (string, error) {
	if clusterType == key.ManagementCluster {
		return "certificate", nil
	}

	secret, err := k8sclient.CoreV1().Secrets(key.Namespace(cluster)).Get(ctx, key.SecretAPICertificates(cluster), metav1.GetOptions{})
	if err != nil {
		return "", microerror.Mask(err)
	}

	if val, ok := secret.Data["token"]; ok && len(val) >= 0 {
		return "token", nil
	}
	return "certificate", nil
}
