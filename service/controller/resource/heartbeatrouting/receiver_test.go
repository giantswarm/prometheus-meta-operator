package heartbeatrouting

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/alertmanager/config"
	commoncfg "github.com/prometheus/common/config"
)

var testReceiver = config.Receiver{
	Name: "test",
	WebhookConfigs: []*config.WebhookConfig{
		&config.WebhookConfig{
			URL: nil,
			HTTPConfig: &commoncfg.HTTPClientConfig{
				BasicAuth: &commoncfg.BasicAuth{
					Password: "pass",
				},
			},
			NotifierConfig: config.NotifierConfig{
				VSendResolved: false,
			},
		},
	},
}

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
						Name: "test",
						WebhookConfigs: []*config.WebhookConfig{
							&config.WebhookConfig{
								URL: nil,
								HTTPConfig: &commoncfg.HTTPClientConfig{
									BasicAuth: &commoncfg.BasicAuth{
										Password: "pass",
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
						Name: "test",
						WebhookConfigs: []*config.WebhookConfig{
							&config.WebhookConfig{
								URL: nil,
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
								URL: nil,
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
			c, updated := ensureReceiver(tc.cfg, testReceiver)
			if updated != tc.expectedUpdate {
				t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
			}

			if len(c.Receivers) != tc.len {
				t.Fatalf("len(c.Receivers) %d, expected %d\n", len(c.Receivers), tc.len)
			}

			if !reflect.DeepEqual(*c.Receivers[tc.index], testReceiver) {
				t.Fatalf("cfg.Receivers[%d] != r\n%v\n", tc.index, cmp.Diff(*c.Receivers[tc.index], testReceiver))
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
						Name: "test",
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
						Name: "test",
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
						Name: "test",
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
						Name: "test",
					},
				},
			},
			expectedUpdate: true,
			len:            0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, updated := removeReceiver(tc.cfg, testReceiver)
			if updated != tc.expectedUpdate {
				t.Fatalf("updated %t, expected %t\n", updated, tc.expectedUpdate)
			}

			if len(c.Receivers) != tc.len {
				t.Fatalf("len(c.Receivers) %d, expected %d\n", len(c.Receivers), tc.len)
			}
		})
	}
}
