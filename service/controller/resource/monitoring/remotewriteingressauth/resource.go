package remotewriteingressauth

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteingressauth"
)

type Config struct {
	K8sClient       k8sclient.Interface
	Logger          micrologger.Logger
	PasswordManager password.Manager
	Installation    string
	Provider        cluster.Provider
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

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.RemoteWriteIngressAuthSecretName,
		Namespace: key.Namespace(cluster),
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func (r *Resource) toSecret(ctx context.Context, v interface{}) (*corev1.Secret, error) {
	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.config.Logger.Debugf(ctx, "looking up for remote write secret")
	username, password, err := remotewriteconfiguration.GetUsernameAndPassword(r.config.K8sClient.K8sClient(), ctx, cluster, r.config.Installation, r.config.Provider)
	if err != nil {
		r.config.Logger.Errorf(ctx, err, "lookup for remote write secret failed")
		return nil, microerror.Mask(err)
	}

	r.config.Logger.Debugf(ctx, "hashing password for the prometheus agent")
	hashedPassword, err := r.config.PasswordManager.Hash([]byte(password))
	if err != nil {
		r.config.Logger.Errorf(ctx, err, "failed to hash the prometheus agent password")
		return nil, microerror.Mask(err)
	}

	r.config.Logger.Debugf(ctx, "hashed password for the prometheus agent")

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			// create authentication string to configure basic auth in nginx
			// see https://github.com/kubernetes/ingress-nginx/blob/main/docs/user-guide/nginx-configuration/annotations.md#authentication
			"auth": []byte(fmt.Sprintf("%s:%s", username, hashedPassword)),
		},
		Type: "Opaque",
	}
	return secret, nil
}

func (r *Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
