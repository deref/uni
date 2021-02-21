package internal

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type PackageMetadata struct {
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Version      string            `json:"version,omitempty"`
	Private      bool              `json:"private"`
	Main         string            `json:"main,omitempty"`
	Dependencies map[string]string `json:"dependencies"`
	Scripts      map[string]string `json:"scripts"`
}

func WritePackageJSON(metadata PackageMetadata, dir string) error {
	metadataPath := path.Join(dir, "package.json")
	return WriteJSON(metadataPath, metadata)
}

func ReadPackageJSON(dir string) (*PackageMetadata, error) {
	metadataPath := path.Join(dir, "package.json")
	bs, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}
	var metadata PackageMetadata
	if err := json.Unmarshal(bs, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}
