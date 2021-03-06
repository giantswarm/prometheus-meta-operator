package rbac

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "rbac"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toClusterRole(v interface{}) (*v1beta1.ClusterRole, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	name := cluster.GetName()

	clusterRole := &v1beta1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: []v1beta1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
					"nodes/metrics",
					"nodes/proxy",
					"services",
					"endpoints",
					"pods",
					"pods/proxy",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"configmaps",
				},
				Verbs: []string{
					"get",
				},
			},
			{
				NonResourceURLs: []string{
					"/metrics",
				},
				Verbs: []string{
					"get",
				},
			},
		},
	}

	return clusterRole, nil
}

func toClusterRoleBinding(v interface{}) (*v1beta1.ClusterRoleBinding, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	name := cluster.GetName()

	clusterRoleBinding := &v1beta1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		RoleRef: v1beta1.RoleRef{
			APIGroup: v1beta1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     name,
		},
		Subjects: []v1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: key.Namespace(cluster),
			},
		},
	}

	return clusterRoleBinding, nil
}

func hasClusterRoleChanged(current, desired *v1beta1.ClusterRole) bool {
	return !reflect.DeepEqual(current.Rules, desired.Rules)
}

func hasClusterRoleBindingChanged(current, desired *v1beta1.ClusterRoleBinding) bool {
	return !reflect.DeepEqual(current.RoleRef, desired.RoleRef) || !reflect.DeepEqual(current.Subjects, desired.Subjects)
}
