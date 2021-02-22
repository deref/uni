package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

type Environment struct {
	Engines []EnvironmentEngine
	Erred   bool
}

type EnvironmentEngine struct {
	Name            string
	ActualVersion   string
	ExpectedVersion string
	Err             error
}

func CheckEngines(repo *Repository) error {
	_, err := AnalyzeEnvironment(repo)
	return err
}

func AnalyzeEnvironment(repo *Repository) (*Environment, error) {
	// Read engine cache.
	engineCache := make(map[string]engineInfo)
	engineCachePath := path.Join(repo.OutDir, "engines.json")
	err := ReadJSON(engineCachePath, &engineCache)
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	// Run engine checks.
	var env Environment
	for engineName, expectedVersion := range repo.Engines {
		info, err := getEngineInfo(engineCache, engineName)
		if err == nil && info.Version != expectedVersion {
			err = fmt.Errorf("unexpected version: %s", info.Version)
		}
		engine := EnvironmentEngine{
			Name:            engineName,
			ActualVersion:   info.Version,
			ExpectedVersion: expectedVersion,
			Err:             err,
		}
		env.Engines = append(env.Engines, engine)
		if err != nil {
			env.Erred = true
		}
	}

	// Write updated engine cache.
	if err := WriteJSON(engineCachePath, engineCache); err != nil {
		return nil, err
	}

	return &env, nil
}

type engineInfo struct {
	Version string    `json:"version"`
	ModTime time.Time `json:"modTime"`
	Size    int64     `json:"size"`
}

func getEngineInfo(cache map[string]engineInfo, name string) (engineInfo, error) {
	args, ok := engineCheckers[name]
	if !ok {
		return engineInfo{}, fmt.Errorf("no engine checker for %q", name)
	}

	// Check binary file.
	binpath, err := exec.LookPath(name)
	if err != nil {
		return engineInfo{}, err
	}
	fileInfo, err := os.Stat(binpath)
	if err != nil {
		return engineInfo{}, err
	}

	// Skip running binary if cached file matches.
	res := engineInfo{
		ModTime: fileInfo.ModTime(),
		Size:    fileInfo.Size(),
	}
	if cached, ok := cache[name]; ok {
		if cached.ModTime.Equal(fileInfo.ModTime()) && cached.Size == fileInfo.Size() {
			res.Version = cached.Version
			return res, nil
		}
	}

	// Run to get version output.
	cmd := exec.Command(binpath, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return engineInfo{}, err
	}
	errstr := string(bytes.TrimSpace(stderr.Bytes()))
	if errstr != "" {
		return engineInfo{}, errors.New(errstr)
	}
	res.Version = string(bytes.TrimSpace(stdout.Bytes()))
	if res.Version == "" {
		return engineInfo{}, errors.New("no version output")
	}

	// Update in-memory engine cache.
	cache[name] = res

	return res, nil
}

var engineCheckers = map[string][]string{
	"node":  {"--version"},
	"npm":   {"--version"},
	"bogus": {},
}

func DumpEnvironment(env *Environment) {
	if len(env.Engines) == 0 {
		fmt.Println("no engine checks")
	} else {
		fmt.Println("engine checks:")
		for _, engine := range env.Engines {
			fmt.Printf("  %s version %s ", engine.Name, engine.ExpectedVersion)
			if engine.Err != nil {
				fmt.Printf("ERROR: %v\n", engine.Err)
			} else {
				fmt.Println("OK")
			}
		}
	}
}
