package internal

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Repository struct {
	ConfigPath   string
	RootDir      string
	OutDir       string
	DistDir      string
	TmpDir       string
	Packages     map[string]*Package
	Dependencies map[string]string
}

type Package struct {
	Name        string
	Public      bool
	Description string
	Entrypoint  string
}

func LoadRepository(searchDir string) (*Repository, error) {
	f, err := openConfigFile(searchDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var repo Repository
	repo.ConfigPath = f.Name()
	repo.RootDir = path.Dir(repo.ConfigPath)
	repo.OutDir = path.Join(repo.RootDir, "out")
	repo.DistDir = path.Join(repo.OutDir, "dist")
	repo.TmpDir = path.Join(repo.OutDir, "tmp")

	dec := yaml.NewDecoder(f)
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	repo.Packages = make(map[string]*Package)
	for packageName, packageConfig := range cfg.Packages {
		repo.Packages[packageName] = &Package{
			Name:        packageName,
			Public:      packageConfig.Public,
			Description: packageConfig.Description,
			Entrypoint:  packageConfig.Entrypoint,
		}
	}

	repo.Dependencies = cfg.Dependencies

	return &repo, nil
}

const configName = "uni.yml"

var ErrNoConfig = fmt.Errorf("cannot find %s config file", configName)

func openConfigFile(searchDir string) (*os.File, error) {
	for {
		configPath := path.Join(searchDir, configName)
		f, err := os.Open(configPath)
		if os.IsNotExist(err) {
			searchDir = path.Dir(searchDir)
			if len(searchDir) <= 1 {
				return nil, ErrNoConfig
			}
			continue
		}
		return f, err
	}
}
