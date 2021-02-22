package internal

import (
	"os"
	"os/exec"
)

type InstallDependenciesOptions struct {
	Frozen bool
}

func InstallDependencies(repo *Repository, opts InstallDependenciesOptions) error {
	metadata := PackageMetadata{
		Name:         "@unirepo/placeholder",
		Private:      true,
		Description:  "GENERATED FILE: DO NOT EDIT! This file is managed by unirepo.",
		Dependencies: repo.Dependencies,
		Scripts: map[string]string{
			"postinstall": "patch-package",
		},
	}
	if err := WritePackageJSON(metadata, repo.RootDir); err != nil {
		return err
	}

	subcommand := "install"
	if opts.Frozen {
		subcommand = "ci"
	}

	npm := exec.Command("npm", subcommand)
	npm.Stdin = os.Stdin
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	return npm.Run()
}
