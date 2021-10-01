package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Fault struct {
	Code int     `yaml:"code"`
	Rate float32 `yaml:"rate"`
}

type Path struct {
	Name         string        `yaml:"name"`
	ResponseTime time.Duration `yaml:"response_time"`
	Faults       []Fault       `yaml:"faults"`
}

type Config struct {
	Paths []Path `yaml:"paths"`
	raw   string
}

func (cfg *Config) Raw() string {
	return cfg.raw
}

func Load(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}

	cfg.raw = string(content)

	return cfg, nil
}
