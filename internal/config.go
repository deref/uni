package internal

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadConfig(filepath string) (*Config, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(bs, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Packages     map[string]PackageConfig
	Dependencies map[string]string
}

type PackageConfig struct {
	Description string
	Entrypoint  string
	// TODO: License, author, etc.
}
