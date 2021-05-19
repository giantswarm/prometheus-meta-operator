package receiver

import (
	"net/url"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	alertmanagerconfig "github.com/giantswarm/prometheus-meta-operator/pkg/alertmanager/config"
	promcommonconfig "github.com/giantswarm/prometheus-meta-operator/pkg/prometheus/common/config"
)

var (
	cluster = &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster",
			Namespace: "cluster-namespace",
		},
	}

	installation = "installation"
	opsgenieKey  = "secret-key"
	u, _         = url.Parse("https://api.opsgenie.com/v2/heartbeats/installation-cluster/ping")
)

func TestEnsureReceiver(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            alertmanagerconfig.Config
		expectedUpdate bool
		len            int
		index          int
	}{
		{
			name: "no update",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "heartbeat_installation_cluster",
						WebhookConfigs: []*alertmanagerconfig.WebhookConfig{
							&alertmanagerconfig.WebhookConfig{
								URL: &alertmanagerconfig.URL{
									URL: u,
								},
								HTTPConfig: &promcommonconfig.HTTPClientConfig{
									BasicAuth: &promcommonconfig.BasicAuth{
										Password: "secret-key",
									},
								},
								NotifierConfig: alertmanagerconfig.NotifierConfig{
									VSendResolved: false,
								},
							},
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
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "heartbeat_installation_cluster",
						WebhookConfigs: []*alertmanagerconfig.WebhookConfig{
							&alertmanagerconfig.WebhookConfig{
								URL: &alertmanagerconfig.URL{
									URL: u,
								},
								HTTPConfig: &promcommonconfig.HTTPClientConfig{
									BasicAuth: &promcommonconfig.BasicAuth{
										Password: "wrong",
									},
								},
								NotifierConfig: alertmanagerconfig.NotifierConfig{
									VSendResolved: false,
								},
							},
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
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "not me",
						WebhookConfigs: []*alertmanagerconfig.WebhookConfig{
							&alertmanagerconfig.WebhookConfig{
								URL: &alertmanagerconfig.URL{
									URL: u,
								},
								HTTPConfig: &promcommonconfig.HTTPClientConfig{
									BasicAuth: &promcommonconfig.BasicAuth{
										Password: "something",
									},
								},
								NotifierConfig: alertmanagerconfig.NotifierConfig{
									VSendResolved: false,
								},
							},
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
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{},
			},
			expectedUpdate: true,
			len:            1,
			index:          0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, updated, err := EnsureCreated(tc.cfg, cluster, installation, opsgenieKey)
			if err != nil {
				t.Fatal(err)
			}

			if updated != tc.expectedUpdate {
				t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
			}

			if len(c.Receivers) != tc.len {
				t.Fatalf("len(c.Receivers) %d, expected %d\n", len(c.Receivers), tc.len)
			}
		})
	}
}

func TestEnsureReceiverWithProxy(t *testing.T) {
	// The HTTPConfig from global config should be merged when creating
	// HTTPConfig for the receiver

	cfg := alertmanagerconfig.Config{
		Global: &alertmanagerconfig.GlobalConfig{
			HTTPConfig: &promcommonconfig.HTTPClientConfig{
				ProxyURL: promcommonconfig.URL{URL: &url.URL{Host: "proxyhost:8080"}},
			},
		},

		Receivers: []*alertmanagerconfig.Receiver{},
	}

	expectedCfg := alertmanagerconfig.Config{
		Global: &alertmanagerconfig.GlobalConfig{
			HTTPConfig: &promcommonconfig.HTTPClientConfig{
				ProxyURL: promcommonconfig.URL{URL: &url.URL{Host: "proxyhost:8080"}},
			},
		},

		Receivers: []*alertmanagerconfig.Receiver{
			{
				Name: "heartbeat_installation_cluster",
				WebhookConfigs: []*alertmanagerconfig.WebhookConfig{
					{
						URL: &alertmanagerconfig.URL{
							URL: u,
						},
						HTTPConfig: &promcommonconfig.HTTPClientConfig{
							BasicAuth: &promcommonconfig.BasicAuth{
								Password: "secret-key",
							},
							ProxyURL: promcommonconfig.URL{URL: &url.URL{Host: "proxyhost:8080"}},
						},
						NotifierConfig: alertmanagerconfig.NotifierConfig{
							VSendResolved: false,
						},
					},
				},
			},
		},
	}

	newCfg, _, err := EnsureCreated(cfg, cluster, installation, opsgenieKey)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(newCfg, expectedCfg) {
		t.Fatalf("\n# expected:\n%+v\n\n# got:\n%+v\n", expectedCfg, newCfg)
	}
}

func TestRemoveReceiver(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            alertmanagerconfig.Config
		expectedUpdate bool
		len            int
	}{
		{
			name: "no update",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "one",
					},
					&alertmanagerconfig.Receiver{
						Name: "two",
					},
				},
			},
			expectedUpdate: false,
			len:            2,
		},
		{
			name: "no update (empty)",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{},
			},
			expectedUpdate: false,
			len:            0,
		},
		{
			name: "remove first",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "heartbeat_installation_cluster",
					},
					&alertmanagerconfig.Receiver{
						Name: "one",
					},
					&alertmanagerconfig.Receiver{
						Name: "two",
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove middle",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "one",
					},
					&alertmanagerconfig.Receiver{
						Name: "heartbeat_installation_cluster",
					},
					&alertmanagerconfig.Receiver{
						Name: "two",
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove last",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "one",
					},
					&alertmanagerconfig.Receiver{
						Name: "two",
					},
					&alertmanagerconfig.Receiver{
						Name: "heartbeat_installation_cluster",
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove (empty)",
			cfg: alertmanagerconfig.Config{
				Receivers: []*alertmanagerconfig.Receiver{
					&alertmanagerconfig.Receiver{
						Name: "heartbeat_installation_cluster",
					},
				},
			},
			expectedUpdate: true,
			len:            0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, updated, err := EnsureDeleted(tc.cfg, cluster, installation, opsgenieKey)
			if err != nil {
				t.Fatal(err)
			}

			if updated != tc.expectedUpdate {
				t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
			}

			if len(c.Receivers) != tc.len {
				t.Fatalf("len(c.Receivers) %d, expected %d\n", len(c.Receivers), tc.len)
			}
		})
	}
}
