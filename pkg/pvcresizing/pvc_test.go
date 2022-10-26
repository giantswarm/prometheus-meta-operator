package pvcresizing

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrometheusVolumeSize(t *testing.T) {

	type args struct {
		value string
	}

	type want struct {
		result string
	}

	expectedLargeValue := "200Gi"
	expectedMediumValue := "100Gi"
	expectedSmallValue := "30Gi"

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"VolumeSizeLarge": {
			reason: "Return volume size large",
			args: args{
				value: "large",
			},
			want: want{
				result: expectedLargeValue,
			},
		},
		"VolumeSizeMedium": {
			reason: "Return volume size medium",
			args: args{
				value: "medium",
			},
			want: want{
				result: expectedMediumValue,
			},
		},
		"VolumeSizeSmall": {
			reason: "Return volume size small",
			args: args{
				value: "small",
			},
			want: want{
				result: expectedSmallValue,
			},
		},
		"VolumeSizeDefault": {
			reason: "Return volume size defaul",
			args: args{
				value: "",
			},
			want: want{
				result: expectedMediumValue,
			},
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			got := PrometheusVolumeSize(tc.args.value)
			if diff := cmp.Diff(tc.want.result, got); diff != "" {
				t.Errorf("\n%s\nExpand(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}
