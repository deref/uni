package internal

import (
	"io/ioutil"
	"os"
)

func EnsureTmp() error {
	return os.MkdirAll(tmpDir, 0755)
}

func TempDir(prefix string) (string, error) {
	return ioutil.TempDir(tmpDir, prefix)
}

func TempFile(prefix string) (*os.File, error) {
	return ioutil.TempFile(tmpDir, "esbuild.meta")
}
