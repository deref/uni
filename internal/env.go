package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
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
	var env Environment
	for engineName, expectedVersion := range repo.Engines {
		actualVersion, err := getEngineVersion(engineName)
		if err == nil && actualVersion != expectedVersion {
			err = fmt.Errorf("unexpected version: %s", actualVersion)
		}
		engine := EnvironmentEngine{
			Name:            engineName,
			ActualVersion:   actualVersion,
			ExpectedVersion: expectedVersion,
			Err:             err,
		}
		env.Engines = append(env.Engines, engine)
		if err != nil {
			env.Erred = true
		}
	}
	return &env, nil
}

func getEngineVersion(name string) (string, error) {
	args, ok := engineCheckers[name]
	if !ok {
		return "", fmt.Errorf("no engine checker for %q", name)
	}
	cmd := exec.Command(name, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	errstr := string(bytes.TrimSpace(stderr.Bytes()))
	if errstr != "" {
		return "", errors.New(errstr)
	}
	version := string(bytes.TrimSpace(stdout.Bytes()))
	if version == "" {
		return "", errors.New("no version output")
	}
	return version, nil
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
