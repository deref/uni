package internal

import "os"

func Clean() error {
	return os.RemoveAll(outDir)
}
