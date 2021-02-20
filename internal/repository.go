package internal

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Repository struct {
	ConfigPath   string
	RootDir      string
	OutDir       string
	TmpDir       string
	Packages     map[string]*Package
	Dependencies map[string]string
}

type Package struct {
	Name        string
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
			Description: packageConfig.Description,
			Entrypoint:  packageConfig.Entrypoint,
		}
	}

	repo.Dependencies = cfg.Dependencies

	return &repo, nil
}

var ErrNoConfig = errors.New("cannot find config file")

func openConfigFile(searchDir string) (*os.File, error) {
	const configName = "uni.yml"
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
