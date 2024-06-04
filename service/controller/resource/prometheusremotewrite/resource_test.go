package prometheusremotewrite

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
)

const (
	name            = "simple-remotewrite"
	namespace       = "default"
	clusterSelector = "giant-cluster"
)

var _ = flag.Bool("update", false, "doing nothing")

func TestEnsurePrometheusRemoteWrite(t *testing.T) {
	r := NewResource()
	rw := remoteWrite(name, namespace, clusterSelector)
	prom := prometheus()
	expectedEmptyPrometheus := expectedPrometheusEmptyRemoteWrite(rw, prom)

	rwAppend := remoteWrite("remotewrite-append", namespace, clusterSelector)
	expectedPrometheusAppend := expectedPrometheusAppend(rwAppend, *expectedEmptyPrometheus.p.DeepCopy())

	rwUpdate := *rw.DeepCopy()
	rwUpdate.Spec.RemoteWrite.URL = "http://my-proxy-url/needs-update"
	expectedPrometheusUpdate := expectedPrometheusUpdate(rwUpdate, *expectedEmptyPrometheus.p.DeepCopy(), 0)

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
			want: want(expectedEmptyPrometheus),
		},
		"AppendingRemoteWriteConfig": {
			reason: "Appending Prometheus remote write config",
			args: args{
				rw: rwAppend,
				p:  *expectedEmptyPrometheus.p.DeepCopy(),
			},
			want: want(expectedPrometheusAppend),
		},
		"UpdateRemoteWriteConfig": {
			reason: "Update current Prometheus remote write config",
			args: args{
				rw: rwUpdate,
				p:  *expectedEmptyPrometheus.p.DeepCopy(),
			},
			want: want(expectedPrometheusUpdate),
		},
		"NoUpdateRemoteWriteConfig": {
			reason: "No update for Prometheus remote write config",
			args: args{
				rw: rw,
				p:  *expectedEmptyPrometheus.p.DeepCopy(),
			},
			want: want{
				p:  expectedEmptyPrometheus.p.DeepCopy(),
				ok: false,
			},
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
	expectedPrometheus := expectedPrometheusEmptyRemoteWrite(rw, prom)
	expectedRemoved := expectedRemoveRemoteWrite(*expectedPrometheus.p.DeepCopy(), 0)

	rwNotFound := remoteWrite("remotewrite-notfound", namespace, clusterSelector)
	expectedNotFound := expectedRemoveRemoteWriteNotExists(*expectedPrometheus.p.DeepCopy())

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
				p:  *expectedPrometheus.p.DeepCopy(),
			},
			want: want(expectedRemoved),
		},
		"RemoteWriteNotFound": {
			reason: "RemoteWrite not found in Prometheus",
			args: args{
				rw: rwNotFound,
				p:  *expectedPrometheus.p.DeepCopy(),
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
				URL:  "https://my-proxy-url",
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
			CommonPrometheusFields: promv1.CommonPrometheusFields{
				Version:  "v1.35.0",
				Replicas: &replicas,
			},
		},
	}
}

func expectedPrometheusEmptyRemoteWrite(rw pmov1alpha1.RemoteWrite, prom promv1.Prometheus) struct {
	p  *promv1.Prometheus
	ok bool
} {
	rw.Spec.RemoteWrite.Name = rw.GetName()

	prom.Spec.RemoteWrite = []promv1.RemoteWriteSpec{rw.Spec.RemoteWrite}
	ok := true

	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: &prom, ok: ok}
}

func expectedPrometheusAppend(rw pmov1alpha1.RemoteWrite, prom promv1.Prometheus) struct {
	p  *promv1.Prometheus
	ok bool
} {
	rw.Spec.RemoteWrite.Name = rw.GetName()

	prom.Spec.RemoteWrite = append(prom.Spec.RemoteWrite, rw.Spec.RemoteWrite)
	ok := true

	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: &prom, ok: ok}
}

func expectedPrometheusUpdate(rw pmov1alpha1.RemoteWrite, prom promv1.Prometheus, rwIndex int) struct {
	p  *promv1.Prometheus
	ok bool
} {
	rw.Spec.RemoteWrite.Name = rw.GetName()

	ok := false
	if !cmp.Equal(rw.Spec.RemoteWrite, prom.Spec.RemoteWrite[rwIndex]) {
		prom.Spec.RemoteWrite[rwIndex] = rw.Spec.RemoteWrite
		ok = true
	}

	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: &prom, ok: ok}
}

func expectedRemoveRemoteWrite(prom promv1.Prometheus, rwIndex int) struct {
	p  *promv1.Prometheus
	ok bool
} {

	b := false
	prom.Spec.RemoteWrite = remove(prom.Spec.RemoteWrite, rwIndex)
	b = true

	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: &prom, ok: b}
}

func expectedRemoveRemoteWriteNotExists(prom promv1.Prometheus) struct {
	p  *promv1.Prometheus
	ok bool
} {

	return struct {
		p  *promv1.Prometheus
		ok bool
	}{p: &prom, ok: false}
}

func NewResource() *Resource {
	r, _ := New(Config{})
	return r
}
