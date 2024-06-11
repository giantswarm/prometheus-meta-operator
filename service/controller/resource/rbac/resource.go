package rbac

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "rbac"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	MimirEnabled bool
}

type Resource struct {
	config Config
}

func New(config Config) (*Resource, error) {
	return &Resource{config}, nil
}

func (r *Resource) Name() string {
	return Name
}

func toClusterRoleBinding(v interface{}) (*v1.ClusterRoleBinding, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	clusterRoleBinding := &v1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: key.Namespace(cluster),
		},
		RoleRef: v1.RoleRef{
			APIGroup: v1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     "prometheus",
		},
		Subjects: []v1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: key.Namespace(cluster),
			},
		},
	}

	return clusterRoleBinding, nil
}

func hasClusterRoleBindingChanged(current, desired *v1.ClusterRoleBinding) bool {
	return !reflect.DeepEqual(current.RoleRef, desired.RoleRef) || !reflect.DeepEqual(current.Subjects, desired.Subjects)
}
