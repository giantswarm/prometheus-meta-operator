package receiver

import (
	"net/url"
	"testing"

	"github.com/prometheus/alertmanager/config"
	commoncfg "github.com/prometheus/common/config"
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
	opsgenieKey  = "secret-key"
	u, _         = url.Parse("https://api.opsgenie.com/v2/heartbeats/installation-cluster/ping")
)

func TestEnsureReceiver(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            config.Config
		expectedUpdate bool
		len            int
		index          int
	}{
		{
			name: "no update",
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "heartbeat_installation_cluster",
						WebhookConfigs: []*config.WebhookConfig{
							&config.WebhookConfig{
								URL: &config.URL{
									URL: u,
								},
								HTTPConfig: &commoncfg.HTTPClientConfig{
									BasicAuth: &commoncfg.BasicAuth{
										Password: "secret-key",
									},
								},
								NotifierConfig: config.NotifierConfig{
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
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "heartbeat_installation_cluster",
						WebhookConfigs: []*config.WebhookConfig{
							&config.WebhookConfig{
								URL: &config.URL{
									URL: u,
								},
								HTTPConfig: &commoncfg.HTTPClientConfig{
									BasicAuth: &commoncfg.BasicAuth{
										Password: "wrong",
									},
								},
								NotifierConfig: config.NotifierConfig{
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
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "not me",
						WebhookConfigs: []*config.WebhookConfig{
							&config.WebhookConfig{
								URL: &config.URL{
									URL: u,
								},
								HTTPConfig: &commoncfg.HTTPClientConfig{
									BasicAuth: &commoncfg.BasicAuth{
										Password: "something",
									},
								},
								NotifierConfig: config.NotifierConfig{
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
			cfg: config.Config{
				Receivers: []*config.Receiver{},
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

func TestRemoveReceiver(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            config.Config
		expectedUpdate bool
		len            int
	}{
		{
			name: "no update",
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "one",
					},
					&config.Receiver{
						Name: "two",
					},
				},
			},
			expectedUpdate: false,
			len:            2,
		},
		{
			name: "no update (empty)",
			cfg: config.Config{
				Receivers: []*config.Receiver{},
			},
			expectedUpdate: false,
			len:            0,
		},
		{
			name: "remove first",
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "heartbeat_installation_cluster",
					},
					&config.Receiver{
						Name: "one",
					},
					&config.Receiver{
						Name: "two",
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove middle",
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "one",
					},
					&config.Receiver{
						Name: "heartbeat_installation_cluster",
					},
					&config.Receiver{
						Name: "two",
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove last",
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
						Name: "one",
					},
					&config.Receiver{
						Name: "two",
					},
					&config.Receiver{
						Name: "heartbeat_installation_cluster",
					},
				},
			},
			expectedUpdate: true,
			len:            2,
		},
		{
			name: "remove (empty)",
			cfg: config.Config{
				Receivers: []*config.Receiver{
					&config.Receiver{
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
