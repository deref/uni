package internal

import (
	"os"
	"os/exec"
)

func Deps(repo *Repository) error {
	metadata := PackageMetadata{
		Private:      true,
		Description:  "GENERATED FILE: DO NOT EDIT! This file is managed by unirepo.",
		Dependencies: repo.Dependencies,
	}
	if err := WritePackageJSON(metadata, repo.RootDir); err != nil {
		return err
	}

	npm := exec.Command("npm", "install")
	npm.Stdin = os.Stdin
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	return npm.Run()
}
