package route

import (
	"testing"

	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	alertmanagerconfig "github.com/giantswarm/prometheus-meta-operator/pkg/alertmanager/config"
)

var (
	cluster = &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster",
			Namespace: "cluster-namespace",
		},
	}

	installation = "installation"
)

func TestEnsureCreated(t *testing.T) {
	var groupInterval, _ = model.ParseDuration("30s")
	var groupWait, _ = model.ParseDuration("5m")
	var repeatInterval, _ = model.ParseDuration("1m")

	testCases := []struct {
		name           string
		cfg            alertmanagerconfig.Config
		expectedUpdate bool
		errorMatcher   func(error) bool
		len            int
	}{
		{
			name: "no update",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "heartbeat_installation_cluster",
							Match: map[string]string{
								"cluster_id":   "cluster",
								"installation": "installation",
								"type":         "heartbeat",
							},
							Continue:       false,
							GroupInterval:  &groupInterval,
							GroupWait:      &groupWait,
							RepeatInterval: &repeatInterval,
						},
					},
				},
			},
			expectedUpdate: false,
			len:            1,
		},
		{
			name: "update",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "heartbeat_installation_cluster",
							Match: map[string]string{
								"cluster_id":   "cluster",
								"installation": "installation",
								"type":         "wrong",
							},
							Continue:       false,
							GroupInterval:  &groupInterval,
							GroupWait:      &groupWait,
							RepeatInterval: &repeatInterval,
						},
					},
				},
			},
			expectedUpdate: true,
			len:            1,
		},
		{
			name: "add",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "not me",
							Match: map[string]string{
								"cluster":      "cluster",
								"installation": "installation",
								"type":         "heartbeat",
							},
							Continue:       false,
							GroupInterval:  &groupInterval,
							GroupWait:      &groupWait,
							RepeatInterval: &repeatInterval,
						},
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "add from 0",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes:   []*alertmanagerconfig.Route{},
				},
			},
			expectedUpdate: true,
			len:            1,
		},
		{
			name: "empty route",
			cfg: alertmanagerconfig.Config{
				Route: nil,
			},
			errorMatcher: IsEmptyRouteError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, updated, err := EnsureCreated(tc.cfg, cluster, installation)
			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
				return
			}

			if err == nil {
				if updated != tc.expectedUpdate {
					t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
				}

				if len(c.Route.Routes) != tc.len {
					t.Fatalf("len(c.Route.Routes) %d, expected %d\n", len(c.Route.Routes), tc.len)
				}
			}
		})
	}
}

func TestEnsureDeleted(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            alertmanagerconfig.Config
		expectedUpdate bool
		errorMatcher   func(error) bool
		len            int
	}{
		{
			name: "no update",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "one",
						},
						&alertmanagerconfig.Route{
							Receiver: "two",
						},
					},
				},
			},
			expectedUpdate: false,
			len:            2,
		},
		{
			name: "no update (empty)",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes:   []*alertmanagerconfig.Route{},
				},
			},
			expectedUpdate: false,
			len:            0,
		},
		{
			name: "remove first",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "heartbeat_installation_cluster",
						},
						&alertmanagerconfig.Route{
							Receiver: "one",
						},
						&alertmanagerconfig.Route{
							Receiver: "two",
						},
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove middle",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "one",
						},
						&alertmanagerconfig.Route{
							Receiver: "heartbeat_installation_cluster",
						},
						&alertmanagerconfig.Route{
							Receiver: "two",
						},
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove last",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "one",
						},
						&alertmanagerconfig.Route{
							Receiver: "two",
						},
						&alertmanagerconfig.Route{
							Receiver: "heartbeat_installation_cluster",
						},
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove (empty)",
			cfg: alertmanagerconfig.Config{
				Route: &alertmanagerconfig.Route{
					Receiver: "root",
					Routes: []*alertmanagerconfig.Route{
						&alertmanagerconfig.Route{
							Receiver: "heartbeat_installation_cluster",
						},
					},
				},
			},
			expectedUpdate: true,
			len:            0,
		},
		{
			name: "empty route",
			cfg: alertmanagerconfig.Config{
				Route: nil,
			},
			expectedUpdate: false,
			len:            0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, updated, err := EnsureDeleted(tc.cfg, cluster, installation)
			if err != nil {
				t.Fatal(err)
			}

			if updated != tc.expectedUpdate {
				t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
			}

			if c.Route != nil && len(c.Route.Routes) != tc.len {
				t.Fatalf("len(c.Route.Routes) %d, expected %d\n", len(c.Route.Routes), tc.len)
			}
		})
	}
}
