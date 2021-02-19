package internal

import "path"

var rootDir = "." // XXX Don't do this!
var outDir = path.Join(rootDir, "out")
var tmpDir = path.Join(outDir, "tmp")

type PackageDeclaration struct {
	Name        string
	Description string
	Entrypoint  string
}
