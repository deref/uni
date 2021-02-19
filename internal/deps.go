package internal

import (
	"os"
	"os/exec"
)

func Deps() error {
	// XXX get from config file.
	dependencies := map[string]string{
		"date-fns": "2.17.0",
	}

	metadata := PackageMetadata{
		Private:      true,
		Description:  "GENERATED FILE: DO NOT EDIT! This file is managed by unirepo.",
		Dependencies: dependencies,
	}
	if err := WritePackageJSON(metadata, rootDir); err != nil {
		return err
	}

	npm := exec.Command("npm", "install")
	npm.Stdin = os.Stdin
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	return npm.Run()
}
