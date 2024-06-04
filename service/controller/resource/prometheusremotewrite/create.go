package prometheusremotewrite

import (
	"context"
	"fmt"
	"net/url"

	"github.com/giantswarm/microerror"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewriteutils"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	if r.mimirEnabled {
		r.logger.Debugf(ctx, "mimir is enabled, deleting")
		return r.EnsureDeleted(ctx, obj)
	}
	r.logger.Debugf(ctx, "ensuring prometheus remoteWrite config")
	{
		remoteWrite, err := remotewriteutils.ToRemoteWrite(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		// fetch current prometheus using the selector provided in remoteWrite resource.
		prometheusList, err := remotewriteutils.FetchPrometheusList(ctx, toResourceWrapper(r), remoteWrite)
		if err != nil {
			return microerror.Mask(err)
		}

		for _, current := range prometheusList.Items {
			err = r.setRemoteWrite(ctx, remoteWrite, current)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		err = r.ensureCleanUp(ctx, remoteWrite, prometheusList.Items)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	r.logger.Debugf(ctx, "ensured prometheus remoteWrite config")

	return nil
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetManagedFields(c.GetManagedFields())
}

func (r *Resource) setRemoteWrite(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, prometheus *promv1.Prometheus) error {
	remoteWriteUrl, err := url.Parse(remoteWrite.Spec.RemoteWrite.URL)
	if err != nil {
		return err
	}
	proxyURL, err := r.proxy(remoteWriteUrl)
	if err != nil {
		return err
	}

	if proxyURL != nil {
		remoteWrite.Spec.RemoteWrite.ProxyURL = proxyURL.String()
	}
	desired, ok := r.ensurePrometheusRemoteWrite(*remoteWrite, *prometheus)
	if !ok {
		r.logger.Debugf(ctx, fmt.Sprintf("no update required for Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
		return nil
	}

	r.logger.Debugf(ctx, fmt.Sprintf("updating Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))

	updateMeta(prometheus, desired)
	_, err = r.prometheusClient.MonitoringV1().
		Prometheuses(prometheus.GetNamespace()).
		Update(ctx, desired, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	err = r.ensureStatusCreated(ctx, remoteWrite, prometheus)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ensureStatusCreated(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, prometheus *promv1.Prometheus) error {
	for _, ref := range remoteWrite.Status.ConfiguredPrometheuses {
		if ref.Name == prometheus.GetName() && ref.Namespace == prometheus.GetNamespace() {
			return nil
		}
	}

	newStatus := corev1.ObjectReference{
		Name:      prometheus.GetName(),
		Namespace: prometheus.GetNamespace(),
	}
	remoteWrite.Status.ConfiguredPrometheuses = append(remoteWrite.Status.ConfiguredPrometheuses, newStatus)

	err := r.k8sClient.CtrlClient().Status().Update(ctx, remoteWrite)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ensureCleanUp(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, prometheuses []*promv1.Prometheus) error {
	// Copy the statuses, because it will be overwritten later on.
	statuses := remoteWrite.Status.ConfiguredPrometheuses

	for _, statusRef := range statuses {
		if !inList(statusRef, prometheuses) {
			p, err := r.prometheusClient.MonitoringV1().
				Prometheuses(statusRef.Namespace).
				Get(ctx, statusRef.Name, metav1.GetOptions{})
			if err != nil && !apierrors.IsNotFound(err) {
				return microerror.Mask(err)
			}

			err = r.unsetRemoteWrite(ctx, remoteWrite, prometheusAndMetadata{
				prometheus: p,
				name:       statusRef.Name,
				namespace:  statusRef.Namespace,
			})
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func inList(o corev1.ObjectReference, list []*promv1.Prometheus) bool {
	for _, item := range list {
		if o.Name == item.GetName() && o.Namespace == item.GetNamespace() {
			return true
		}
	}

	return false
}
