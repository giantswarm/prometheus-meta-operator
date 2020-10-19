package route

import (
	"testing"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	var one, _ = model.ParseDuration("1s")
	var fifteen, _ = model.ParseDuration("15s")

	testCases := []struct {
		name           string
		cfg            config.Config
		expectedUpdate bool
		errorMatcher   func(error) bool
		len            int
	}{
		{
			name: "no update",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "heartbeat_installation_cluster",
							Match: map[string]string{
								"cluster_id":   "cluster",
								"installation": "installation",
								"type":         "heartbeat",
							},
							Continue:       false,
							GroupInterval:  &one,
							GroupWait:      &one,
							RepeatInterval: &fifteen,
						},
					},
				},
			},
			expectedUpdate: false,
			len:            1,
		},
		{
			name: "update",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "heartbeat_installation_cluster",
							Match: map[string]string{
								"cluster_id":   "cluster",
								"installation": "installation",
								"type":         "wrong",
							},
							Continue:       false,
							GroupInterval:  &one,
							GroupWait:      &one,
							RepeatInterval: &fifteen,
						},
					},
				},
			},
			expectedUpdate: true,
			len:            1,
		},
		{
			name: "add",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "not me",
							Match: map[string]string{
								"cluster":      "cluster",
								"installation": "installation",
								"type":         "heartbeat",
							},
							Continue:       false,
							GroupInterval:  &one,
							GroupWait:      &one,
							RepeatInterval: &fifteen,
						},
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "add from 0",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes:   []*config.Route{},
				},
			},
			expectedUpdate: true,
			len:            1,
		},
		{
			name: "empty route",
			cfg: config.Config{
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
		cfg            config.Config
		expectedUpdate bool
		errorMatcher   func(error) bool
		len            int
	}{
		{
			name: "no update",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "one",
						},
						&config.Route{
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
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes:   []*config.Route{},
				},
			},
			expectedUpdate: false,
			len:            0,
		},
		{
			name: "remove first",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "heartbeat_installation_cluster",
						},
						&config.Route{
							Receiver: "one",
						},
						&config.Route{
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
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "one",
						},
						&config.Route{
							Receiver: "heartbeat_installation_cluster",
						},
						&config.Route{
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
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "one",
						},
						&config.Route{
							Receiver: "two",
						},
						&config.Route{
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
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
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
			cfg: config.Config{
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
