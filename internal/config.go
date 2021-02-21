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
	// TODO: Do not load these in as maps, since duplicate keys are
	// not decectable this way!
	Packages     map[string]PackageConfig
	Dependencies map[string]string
}

type PackageConfig struct {
	Public      bool
	Description string
	Entrypoint  string
	// TODO: License, author, etc.
}
