package internal

import (
	"fmt"
	"io"
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
	Engines      map[string]string
	Packages     map[string]*Package
	Dependencies map[string]*Dependency
}

type Dependency struct {
	Name    string
	Version string
}

type Package struct {
	Name        string
	Public      bool
	Description string
	Index       string
	Executables map[string]*Executable
}

type Executable struct {
	Name       string
	Entrypoint string
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
	dec.SetStrict(true)
	var cfg Config
	err = dec.Decode(&cfg)
	if err == io.EOF {
		// TODO: Should entirely empty config files be allowed? Probably not.
		err = nil
	}
	if err != nil {
		return nil, err
	}

	repo.Engines = make(map[string]string)
	for engineName, engineVersion := range cfg.Engines {
		repo.Engines[engineName] = engineVersion
	}

	repo.Packages = make(map[string]*Package)
	for packageName, packageConfig := range cfg.Packages {
		pkg := &Package{
			Name:        packageName,
			Public:      packageConfig.Public,
			Description: packageConfig.Description,
			Index:       packageConfig.Index,
		}
		pkg.Executables = make(map[string]*Executable)
		for executableName, executableEntrypoint := range packageConfig.Executables {
			pkg.Executables[executableName] = &Executable{
				Name:       executableName,
				Entrypoint: executableEntrypoint,
			}
		}
		repo.Packages[packageName] = pkg
	}

	repo.Dependencies = make(map[string]*Dependency)
	addDependency := func(name, version string) {
		repo.Dependencies[name] = &Dependency{
			Name:    name,
			Version: version,
		}
	}
	for dependencyName, dependencyVersion := range cfg.Dependencies {
		addDependency(dependencyName, dependencyVersion)
	}
	for dependencyName, dependencyVersion := range requiredDependencies {
		if _, ok := repo.Dependencies[dependencyName]; !ok {
			addDependency(dependencyName, dependencyVersion)
		}
	}

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
