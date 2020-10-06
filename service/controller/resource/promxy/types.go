package promxy

import (
	"gopkg.in/yaml.v2"
)

type Promxy struct {
	Promxy PromxyConfig `yaml:"promxy"`
}

// PromxyConfig is the configuration for Promxy itself
type PromxyConfig struct {
	// Config for each of the server groups promxy is configured to aggregate
	ServerGroups []*ServerGroup `yaml:"server_groups"`
}

func (p *PromxyConfig) Contains(group ServerGroup) bool {
	for _, val := range p.ServerGroups {
		if val.PathPrefix == group.PathPrefix {
			return true
		}
	}
	return false
}

func (p *PromxyConfig) Add(group ServerGroup) {
	p.ServerGroups = append(p.ServerGroups, &group)
}

func (p *PromxyConfig) Remove(group ServerGroup) {
	var index int
	for key, val := range p.ServerGroups {
		if val.PathPrefix == group.PathPrefix {
			index = key
		}
	}
	p.ServerGroups = append(p.ServerGroups[:index], p.ServerGroups[index+1:]...)
}

func Serialize(config Promxy) (string, error) {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func Deserialize(content string) (Promxy, error) {
	config := Promxy{}
	err := yaml.Unmarshal([]byte(content), &config)
	return config, err
}
