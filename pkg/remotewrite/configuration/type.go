package configuration

import (
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

type RemoteWriteConfig struct {
	// We keep this until we are sure this can be removed
	GlobalConfig          GlobalConfig          `yaml:"global,omitempty" json:"global,omitempty"`
	PrometheusAgentConfig PrometheusAgentConfig `yaml:"prometheus-agent,omitempty" json:"prometheus-agent,omitempty"`
}

type GlobalConfig struct {
	RemoteWrite    []RemoteWrite     `yaml:"remoteWrite,omitempty" json:"remoteWrite,omitempty"`
	ExternalLabels map[string]string `yaml:"externalLabels,omitempty" json:"externalLabels,omitempty"`
}

type PrometheusAgentConfig struct {
	RemoteWrite    []RemoteWrite        `yaml:"remoteWrite,omitempty" json:"remoteWrite,omitempty"`
	ExternalLabels map[string]string    `yaml:"externalLabels,omitempty" json:"externalLabels,omitempty"`
	Image          PrometheusAgentImage `yaml:"image,omitempty" json:"image,omitempty"`
	Shards         int                  `yaml:"shards,omitempty" json:"shards,omitempty"`
	Version        string               `yaml:"version,omitempty" json:"version,omitempty"`
}

type PrometheusAgentImage struct {
	Tag string `yaml:"tag" json:"tag"`
}

type RemoteWrite struct {
	Name        string             `yaml:"name" json:"name"`
	Password    string             `yaml:"password" json:"password"`
	Username    string             `yaml:"username" json:"username"`
	URL         string             `yaml:"url" json:"url"`
	QueueConfig promv1.QueueConfig `yaml:"queueConfig" json:"queueConfig"`
	TLSConfig   promv1.TLSConfig   `yaml:"tlsConfig" json:"tlsConfig"`
}
