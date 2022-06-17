package prometheusremotewrite

import (
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
)

const (
	name            = "simple-remotewrite"
	namespace       = "default"
	clusterSelector = "giant-cluster"
)

func TestToRemoteWrite(t *testing.T) {

	type args struct {
		obj interface{}
	}

	type want struct {
		rw  *pmov1alpha1.RemoteWrite
		err error
	}

	successObj := pmov1alpha1.RemoteWrite{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: pmov1alpha1.RemoteWriteSpec{},
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ConvertSuccess": {
			reason: "Convert an object to RemoteWrite",
			args: args{
				obj: &pmov1alpha1.RemoteWrite{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: pmov1alpha1.RemoteWriteSpec{},
				},
			},
			want: want{
				err: nil,
				rw:  &successObj,
			},
		},
		"ConvertFailed": {
			reason: "Convert an object to RemoteWrite Failed",
			args: args{
				obj: promv1.Prometheus{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: promv1.PrometheusSpec{},
				},
			},
			want: want{
				err: microerror.Maskf(wrongTypeError, "'%T' is not a 'pmov1alpha1.RemoteWrite'", promv1.Prometheus{}),
				rw:  nil,
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			got, err := ToRemoteWrite(tc.args.obj)
			if diff := cmp.Diff(tc.want.err, err, unittest.EquateErrors()); diff != "" {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%s", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.rw, got); diff != "" {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}

func TestEnsurePrometheusRemoteWrite(t *testing.T) {

	r := NewResource()
	rw := remoteWrite(name, namespace, clusterSelector)
	prom := prometheus()
	expectedPrometheus := wrapperPrometheus(rw, prom)

	rwAppend := remoteWrite("remotewrite-append", namespace, clusterSelector)
	expectedPrometheusAppend := wrapperPrometheus(rwAppend, *expectedPrometheus.p)

	rwUpdate := rw
	rwUpdate.Spec.RemoteWrite.URL = "http://my-fancy-url/needs-update"
	expectedPrometheusUpdate := wrapperPrometheus(rwUpdate, *expectedPrometheus.p)

	type args struct {
		rw pmov1alpha1.RemoteWrite
		p  promv1.Prometheus
	}

	type want struct {
		p  *promv1.Prometheus
		ok bool
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"AddingRemoteWriteConfigOnEmptyRemoteWrite": {
			reason: "Updating empty Prometheus remote write config",
			args: args{
				rw: rw,
				p:  prom,
			},
			want: want(expectedPrometheus),
		},
		"AppendingRemoteWriteConfig": {
			reason: "Appending Prometheus remote write config",
			args: args{
				rw: rwAppend,
				p:  *expectedPrometheus.p,
			},
			want: want(expectedPrometheusAppend),
		},
		"UpdateRemoteWriteConfig": {
			reason: "Update current Prometheus remote write config",
			args: args{
				rw: rwUpdate,
				p:  *expectedPrometheus.p,
			},
			want: want(expectedPrometheusUpdate),
		},
		"NoUpdateRemoteWriteConfig": {
			reason: "No update for Prometheus remote write config",
			args: args{
				rw: rw,
				p:  *expectedPrometheus.p,
			},
			want: want(wrapperPrometheus(rw, *expectedPrometheus.p)),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			got, ok := r.ensurePrometheusRemoteWrite(tc.args.rw, tc.args.p)
			if tc.want.ok != ok {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%v", tc.reason, ok)
			}
			if diff := cmp.Diff(tc.want.p, got); diff != "" {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}

func TestRemovePrometheusRemoteWrite(t *testing.T) {

	rw := remoteWrite(name, namespace, clusterSelector)
	prom := prometheus()
	expectedPrometheus := wrapperPrometheus(rw, prom)
	expectedRemoved := wrapperRemoveRemoteWrite(rw, *expectedPrometheus.p)

	rwNotFound := remoteWrite("remotewrite-notfound", namespace, clusterSelector)
	expectedNotFound := wrapperRemoveRemoteWrite(rwNotFound, *expectedPrometheus.p)

	type args struct {
		rw pmov1alpha1.RemoteWrite
		p  promv1.Prometheus
	}

	type want struct {
		p  *promv1.Prometheus
		ok bool
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"RemoteWriteRemoved": {
			reason: "RemoteWrite removed from Prometheus",
			args: args{
				rw: rw,
				p:  *expectedPrometheus.p,
			},
			want: want(expectedRemoved),
		},
		"RemoteWriteNotFound": {
			reason: "RemoteWrite not found in Prometheus",
			args: args{
				rw: rwNotFound,
				p:  *expectedPrometheus.p,
			},
			want: want(expectedNotFound),
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			got, ok := removePrometheusRemoteWrite(tc.args.rw, tc.args.p)
			if tc.want.ok != ok {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%v", tc.reason, ok)
			}
			if diff := cmp.Diff(tc.want.p, got); diff != "" {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}

func remoteWrite(name, namespace, clusterSelector string) pmov1alpha1.RemoteWrite {

	return pmov1alpha1.RemoteWrite{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: pmov1alpha1.RemoteWriteSpec{
			ClusterSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"cluster": clusterSelector},
			},
			RemoteWrite: promv1.RemoteWriteSpec{
				URL:  "https://my-fancy-url",
				Name: "test",
				BasicAuth: &promv1.BasicAuth{
					Username: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: name,
						},
						Key: "username",
					},
					Password: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: name,
						},
						Key: "password",
					},
				},
			},
		},
	}
}

func prometheus() promv1.Prometheus {
	var replicas int32 = 1
	return promv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-prometheus",
			Namespace: "default",
		},
		Spec: promv1.PrometheusSpec{
			Version:  "v1.35.0",
			Replicas: &replicas,
		},
	}
}

func wrapperPrometheus(rw pmov1alpha1.RemoteWrite, prom promv1.Prometheus) struct {
	p  *promv1.Prometheus
	ok bool
} {
	r := NewResource()
	p, ok := r.ensurePrometheusRemoteWrite(rw, prom)
	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: p, ok: ok}
}

func wrapperRemoveRemoteWrite(r pmov1alpha1.RemoteWrite, prom promv1.Prometheus) struct {
	p  *promv1.Prometheus
	ok bool
} {

	p, ok := removePrometheusRemoteWrite(r, prom)
	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: p, ok: ok}
}

func NewResource() *Resource {
	config := Config{
		HTTPProxy:  "http://proxy-url",
		HTTPSProxy: "",
		NoProxy:    "http://no-proxy",
	}
	r, _ := New(config)
	return r
}
