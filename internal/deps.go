package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type PackageMetadata struct {
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Private      bool              `json:"private,omitempty"`
	Dependencies map[string]string `json:"dependencies"`
}

func Deps() error {
	metadata := PackageMetadata{
		Private:      true,
		Description:  "GENERATED FILE: DO NOT EDIT! This file is managed by unirepo.",
		Dependencies: make(map[string]string),
	}
	bs, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataPath := "package.json" // XXX absolute path
	if err := ioutil.WriteFile(metadataPath, bs, 0644); err != nil {
		return err
	}

	npm := exec.Command("npm", "install")
	npm.Stdin = os.Stdin
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	return npm.Run()
}
