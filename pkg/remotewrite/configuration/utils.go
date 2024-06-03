package configuration

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const remoteWriteEndpointTemplateURL = "https://%s/%s/api/v1/write"

func GetUsernameAndPassword(client kubernetes.Interface, ctx context.Context, cluster v1.Object, installation string, provider cluster.Provider) (string, string, error) {
	secretName := key.RemoteWriteSecretName(cluster)
	secretNamespace := key.GetClusterAppsNamespace(cluster, installation, provider)

	remoteWriteSecret, err := client.CoreV1().Secrets(secretNamespace).Get(ctx, secretName, v1.GetOptions{})
	if err != nil {
		return "", "", microerror.Mask(err)
	}

	username, password, err := extractUsernameAndPasswordFromSecret(remoteWriteSecret)
	if err != nil {
		return "", "", microerror.Mask(err)
	}
	return username, password, nil
}

func extractUsernameAndPasswordFromSecret(secret *corev1.Secret) (string, string, error) {
	data := secret.Data["values"]
	remoteWriteValues := RemoteWriteConfig{}
	err := yaml.Unmarshal(data, &remoteWriteValues)
	if err != nil {
		return "", "", microerror.Mask(err)
	}

	if len(remoteWriteValues.GlobalConfig.RemoteWrite) == 0 && len(remoteWriteValues.PrometheusAgentConfig.RemoteWrite) == 0 {
		// skipping
		return "", "", secretNotFound
	}

	for _, rw := range remoteWriteValues.GlobalConfig.RemoteWrite {
		if rw.Name == key.PrometheusMetaOperatorRemoteWriteName {
			return rw.Username, rw.Password, nil
		}
	}

	for _, rw := range remoteWriteValues.PrometheusAgentConfig.RemoteWrite {
		if rw.Name == key.PrometheusMetaOperatorRemoteWriteName {
			return rw.Username, rw.Password, nil
		}
	}

	return "", "", remoteWriteNotFound
}

func DefaultRemoteWrite(clusterID string, baseDomain string, password string, insecureCA bool) RemoteWrite {
	url := fmt.Sprintf(remoteWriteEndpointTemplateURL, baseDomain, clusterID)
	return RemoteWrite{
		Name:          key.PrometheusMetaOperatorRemoteWriteName,
		URL:           url,
		Username:      clusterID,
		Password:      password,
		RemoteTimeout: "60s",
		QueueConfig: promv1.QueueConfig{
			Capacity:          30000,
			MaxSamplesPerSend: 150000,
			MaxShards:         10,
		},
		TLSConfig: promv1.TLSConfig{
			SafeTLSConfig: promv1.SafeTLSConfig{
				InsecureSkipVerify: &insecureCA,
			},
		},
	}
}
