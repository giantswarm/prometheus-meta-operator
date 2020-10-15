package heartbeatrouting

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/alertmanager/config"
)

var testRoute = config.Route{
	Receiver: "test",
	Match: map[string]string{
		"cluster":      "cluster",
		"installation": "installation",
		"type":         "heartbeat",
	},
	Continue:       false,
	GroupInterval:  &one,
	GroupWait:      &one,
	RepeatInterval: &fifteen,
}

func TestEnsureRoute(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            config.Config
		expectedUpdate bool
		errorMatcher   func(error) bool
		len            int
		index          int
	}{
		{
			name: "no update",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "test",
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
			expectedUpdate: false,
			len:            1,
			index:          0,
		},
		{
			name: "update",
			cfg: config.Config{
				Route: &config.Route{
					Receiver: "root",
					Routes: []*config.Route{
						&config.Route{
							Receiver: "test",
							Match: map[string]string{
								"cluster":      "cluster",
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
			index:          0,
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
			index:          1,
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
			index:          0,
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
			c, updated, err := ensureRoute(tc.cfg, testRoute)
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

				if !reflect.DeepEqual(*c.Route.Routes[tc.index], testRoute) {
					t.Fatalf("cfg.Route.Routes[%d] != r\n%v\n", tc.index, cmp.Diff(*c.Route.Routes[tc.index], testRoute))
				}
			}
		})
	}
}

func TestRemoveRoute(t *testing.T) {
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
							Receiver: "test",
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
							Receiver: "test",
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
							Receiver: "test",
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
							Receiver: "test",
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
			c, updated := removeRoute(tc.cfg, testRoute)
			if updated != tc.expectedUpdate {
				t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
			}

			if c.Route != nil && len(c.Route.Routes) != tc.len {
				t.Fatalf("len(c.Route.Routes) %d, expected %d\n", len(c.Route.Routes), tc.len)
			}
		})
	}
}
