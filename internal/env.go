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
	OK      bool
}

type EnvironmentEngine struct {
	Name            string
	ActualVersion   string
	ExpectedVersion string
	OK              bool
}

func CheckEngines(repo *Repository) error {
	env, err := AnalyzeEnvironment(repo)
	if err != nil {
		return err
	}
	if env.OK {
		return nil
	}
	for _, engine := range env.Engines {
		if !engine.OK {
			return fmt.Errorf("engine error: expected %s version %s, but have %s", engine.Name, engine.ExpectedVersion, engine.ActualVersion)
		}
	}
	panic("unreachable")
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
		Warnf("failed to read engine cache: %v", err)
	}

	// Run engine checks.
	var env Environment
	env.OK = true
	for engineName, expectedVersion := range repo.Engines {
		info, err := getEngineInfo(engineCache, engineName)
		if err != nil {
			return nil, fmt.Errorf("error checking %s: %w", engineName, err)
		}
		engine := EnvironmentEngine{
			Name:            engineName,
			ActualVersion:   info.Version,
			ExpectedVersion: expectedVersion,
			OK:              info.Version == expectedVersion,
		}
		env.Engines = append(env.Engines, engine)
		env.OK = env.OK && engine.OK
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

// side-effect: updates cache.
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
			if engine.OK {
				fmt.Println("OK")
			} else {
				fmt.Printf("ERROR: actual version %s\n", engine.ActualVersion)
			}
		}
	}
}
