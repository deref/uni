package internal

import "os"

func Clean(repo *Repository) error {
	return os.RemoveAll(repo.OutDir)
}
