package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/natefinch/atomic"
	"golang.org/x/sync/errgroup"
)

func ReadJSON(filename string, data interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err1 := f.Close(); err == nil {
			err = err1
		}
	}()
	dec := json.NewDecoder(f)
	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("error decoding %q: %w", filename, err)
	}
	return
}

func WriteJSON(filename string, data interface{}) (err error) {
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}

	r, w := io.Pipe()

	var eg errgroup.Group
	eg.Go(func() error {
		defer w.Close()
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	})

	eg.Go(func() error {
		return atomic.WriteFile(filename, r)
	})

	return eg.Wait()
}
