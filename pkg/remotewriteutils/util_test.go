package remotewriteutils

import (
	"flag"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

const (
	name      = "simple-remotewrite"
	namespace = "default"
)

var _ = flag.Bool("update", false, "doing nothing")

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
