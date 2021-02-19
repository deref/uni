package internal

import (
	"path"
)

type PackageMetadata struct {
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Private      bool              `json:"private,omitempty"`
	Main         string            `json:"main,omitempty"`
	Dependencies map[string]string `json:"dependencies"`
}

func WritePackageJSON(metadata PackageMetadata, dir string) error {
	filepath := path.Join(dir, "package.json")
	return WriteJSON(filepath, metadata)
}
