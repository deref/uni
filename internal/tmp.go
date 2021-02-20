package internal

import (
	"io/ioutil"
	"os"
)

func EnsureTmp(repo *Repository) error {
	return os.MkdirAll(repo.TmpDir, 0755)
}

func TempDir(repo *Repository, prefix string) (string, error) {
	return ioutil.TempDir(repo.TmpDir, prefix)
}

func TempFile(repo *Repository, prefix string) (*os.File, error) {
	return ioutil.TempFile(repo.TmpDir, "esbuild.meta")
}
