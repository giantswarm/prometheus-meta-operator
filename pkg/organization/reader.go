package organization

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	DefaultOrganization string = "giantswarm"
	OrganizationLabel   string = "giantswarm.io/organization"
)

type Reader interface {
	Read(ctx context.Context, cluster metav1.Object) (string, error)
}

type NamespaceReader struct {
	client       kubernetes.Interface
	installation string
	provider     cluster.Provider
}

func NewNamespaceReader(client kubernetes.Interface, installation string, provider cluster.Provider) Reader {
	return NamespaceReader{client, installation, provider}
}

func (r NamespaceReader) Read(ctx context.Context, cluster metav1.Object) (string, error) {
	// Vintage MC
	if key.IsManagementCluster(r.installation, cluster) && !key.IsCAPIManagementCluster(r.provider) {
		return DefaultOrganization, nil
	}

	// For the rest, we extract the organization name from the namespace labels
	namespace, err := r.client.CoreV1().Namespaces().Get(ctx, cluster.GetNamespace(), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if organization, ok := namespace.Labels[OrganizationLabel]; ok {
		return organization, nil
	}
	return "", errors.New("cluster namespace missing organization label")
}
