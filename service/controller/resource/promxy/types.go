package promxy

import (
	"fmt"
	"time"

	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/pkg/relabel"
)

// Promxy defines the content of promxy configuration file.
// This structure is copied over from godoc.org/github.com/jacksontj/promxy/pkg/config#Config as Promxy
// is relying on prometheus 1.8 and we are using prometheus 2.20 with implies some changes in prometheus types
// and makes the defautl promxy struct incompatible with our codebase.
// TODO, move to promxy type once https://github.com/jacksontj/promxy/issues/352 is resolved
type Promxy struct {
	// Prometheus configs that includes configurations for
	// recording rules, alerting rules, etc.
	PromConfig config.Config `yaml:",inline"`
	// Promxy specific configuration
	Promxy PromxyConfig `yaml:"promxy"`
}

// PromxyConfig is the configuration for Promxy itself
type PromxyConfig struct {
	// Config for each of the server groups promxy is configured to aggregate
	ServerGroups []*ServerGroup `yaml:"server_groups"`
}

func (p *PromxyConfig) contains(group *ServerGroup) bool {
	for _, val := range p.ServerGroups {
		if val.PathPrefix == group.PathPrefix {
			return true
		}
	}
	return false
}

func (p *PromxyConfig) add(group *ServerGroup) {
	p.ServerGroups = append(p.ServerGroups, group)
}

func (p *PromxyConfig) remove(group *ServerGroup) {
	var index int
	for key, val := range p.ServerGroups {
		if val.PathPrefix == group.PathPrefix {
			index = key
		}
	}
	p.ServerGroups = append(p.ServerGroups[:index], p.ServerGroups[index+1:]...)
}

// ServerGroup is the configuration for a ServerGroup that promxy will talk to.
// This is where the vast majority of options exist.
type ServerGroup struct {
	// RemoteRead directs promxy to load RAW data (meaning matrix selectors such as `foo[1h]`)
	// through the RemoteRead API on prom.
	// Pros:
	//  - StaleNaNs work
	//  - ~2x faster (in my local testing, more so if you are using default JSON marshaler in prom)
	//
	// Cons:
	//  - proto marshaling prom side doesn't stream, so the data being sent
	//      over the wire will be 2x its size in memory on the remote prom host.
	//  - "experimental" API (according to docs) -- meaning this might break
	//      without much (if any) warning
	//
	// Upstream prom added a StaleNan to determine if a given timeseries has gone
	// NaN -- the problem being that for range vectors they filter out all "stale" samples
	// meaning that it isn't possible to get a "raw" dump of data through the query/query_range v1 API
	// The only option that exists in reality is the "remote read" API -- which suffers
	// from the same memory-balooning problems that the HTTP+JSON API originally had.
	// It has **less** of a problem (its 2x memory instead of 14x) so it is a viable option.
	RemoteRead bool `yaml:"remote_read"`
	// RemoteReadPath sets the remote read path for the hosts in this servergroup
	RemoteReadPath string `yaml:"remote_read_path"`
	// HTTP client config for promxy to use when connecting to the various server_groups
	// this is the same config as prometheus
	HTTPConfig HTTPClientConfig `yaml:"http_client"`
	// Scheme defines how promxy talks to this server group (http, https, etc.)
	Scheme string `yaml:"scheme"`
	// Labels is a set of labels that will be added to all metrics retrieved
	// from this server group
	Labels model.LabelSet `yaml:"labels"`
	// RelabelConfigs are similar in function and identical in configuration as prometheus'
	// relabel config for scrape jobs. The difference here being that the source labels
	// you can pull from are from the downstream servergroup target and the labels you are
	// relabeling are that of the timeseries being returned. This allows you to mutate the
	// labelsets returned by that target at runtime.
	// To further illustrate the difference we'll look at an example:
	//
	//      relabel_configs:
	//    - source_labels: [__meta_consul_tags]
	//      regex: '.*,prod,.*'
	//      action: keep
	//    - source_labels: [__meta_consul_dc]
	//      regex: '.+'
	//      action: replace
	//      target_label: datacenter
	//
	// If we saw this in a scrape-config we would expect:
	//   (1) the scrape would only target hosts with a prod consul label
	//   (2) it would add a label to all returned series of datacenter with the value set to whatever the value of __meat_consul_dc was.
	//
	// If we saw this same config in promxy (pointing at prometheus hosts instead of some exporter), we'd expect a similar behavior:
	//   (1) only targets with the prod consul label would be included in the servergroup
	//   (2) it would add a label to all returned series of this servergroup of datacenter with the value set to whatever the value of __meat_consul_dc was.
	//
	// So in reality its "the same", the difference is in prometheus these apply to the labels/targets of a scrape job,
	// in promxy they apply to the prometheus hosts in the servergroup - but the behavior is the same.
	RelabelConfigs []*relabel.Config `yaml:"relabel_configs,omitempty"`
	// Hosts is a set of discovery.Config options that allow promxy to discover
	// all hosts in the server_group
	KubernetesSDConfigs []*kubernetes.SDConfig `yaml:"kubernetes_sd_configs,omitempty"`
	// PathPrefix to prepend to all queries to hosts in this servergroup
	PathPrefix string `yaml:"path_prefix"`
	// QueryParams are a map of query params to add to all HTTP calls made to this downstream
	// the main use-case for this is to add `nocache=1` to VictoriaMetrics downstreams
	// (see https://github.com/jacksontj/promxy/issues/202)
	QueryParams map[string]string `yaml:"query_params,omitempty"`
	// TODO cache this as a model.Time after unmarshal
	// AntiAffinity defines how large of a gap in the timeseries will cause promxy
	// to merge series from 2 hosts in a server_group. This required for a couple reasons
	// (1) Promxy cannot make assumptions on downstream clock-drift and
	// (2) two prometheus hosts scraping the same target may have different times
	// #2 is caused by prometheus storing the time of the scrape as the time the scrape **starts**.
	// in practice this is actually quite frequent as there are a variety of situations that
	// cause variable scrape completion time (slow exporter, serial exporter, network latency, etc.)
	// any one of these can cause the resulting data in prometheus to have the same time but in reality
	// come from different points in time. Best practice for this value is to set it to your scrape interval
	AntiAffinity time.Duration `yaml:"anti_affinity,omitempty"`

	// IgnoreError will hide all errors from this given servergroup effectively making
	// the responses from this servergroup "not required" for the result.
	// Note: this allows you to make the tradeoff between availability of queries and consistency of results
	IgnoreError bool `yaml:"ignore_error,omitempty"`

	// RelativeTimeRangeConfig defines a relative time range that this servergroup will respond to
	// An example use-case would be if a specific servergroup was long-term storage, it might only
	// have data 3d old and retain 90d of data.
	*RelativeTimeRangeConfig `yaml:"relative_time_range,omitempty"`

	// AbsoluteTimeRangeConfig defines an absolute time range that this servergroup will respond to
	// An example use-case would be if a specific servergroup was was "deprecated" and wasn't getting
	// any new data after a specific given point in time
	*AbsoluteTimeRangeConfig `yaml:"absolute_time_range,omitempty"`
}

// GetScheme returns the scheme for this servergroup
func (c *ServerGroup) GetScheme() string {
	return c.Scheme
}

// GetAntiAffinity returns the AntiAffinity time for this servergroup
func (c *ServerGroup) GetAntiAffinity() model.Time {
	return model.TimeFromUnix(int64((c.AntiAffinity).Seconds()))
}

// HTTPClientConfig extends prometheus' HTTPClientConfig
type HTTPClientConfig struct {
	DialTimeout time.Duration                `yaml:"dial_timeout"`
	HTTPConfig  config_util.HTTPClientConfig `yaml:",inline"`
}

// RelativeTimeRangeConfig configures durations relative from "now" to define
// a servergroup's time range
type RelativeTimeRangeConfig struct {
	Start *time.Duration `yaml:"start"`
	End   *time.Duration `yaml:"end"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (tr *RelativeTimeRangeConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain RelativeTimeRangeConfig
	if err := unmarshal((*plain)(tr)); err != nil {
		return err
	}

	return tr.validate()
}

func (tr *RelativeTimeRangeConfig) validate() error {
	if tr.End != nil && tr.Start != nil && *tr.End < *tr.Start {
		return fmt.Errorf("RelativeTimeRangeConfig: End must be after start")
	}
	return nil
}

// AbsoluteTimeRangeConfig contains absolute times to define a servergroup's time range
type AbsoluteTimeRangeConfig struct {
	Start time.Time `yaml:"start"`
	End   time.Time `yaml:"end"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (tr *AbsoluteTimeRangeConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain AbsoluteTimeRangeConfig
	if err := unmarshal((*plain)(tr)); err != nil {
		return err
	}

	return tr.validate()
}

func (tr *AbsoluteTimeRangeConfig) validate() error {
	if !tr.Start.IsZero() && !tr.End.IsZero() && tr.End.Before(tr.Start) {
		return fmt.Errorf("AbsoluteTimeRangeConfig: End must be after start")
	}
	return nil
}
