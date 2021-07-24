package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
)

type buildAndWatch struct {
	Repository    *Repository
	Package       *Package
	Esbuild       api.BuildOptions // XXX smaller option set.
	Types         bool
	Watch         bool
	CreateProcess func() process
}

type process interface {
	Start() error
	Wait() error
	Kill() error
}

func (opts buildAndWatch) Run() error {
	if opts.Watch && opts.Types {
		return errors.New("cannot build types with watch")
	}

	repo := opts.Repository

	plugins := append([]api.Plugin{}, opts.Esbuild.Plugins...)

	var watcher *fsnotify.Watcher
	if opts.Watch {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		watchPlugin := api.Plugin{
			Name: "monoclean:watch",
			Setup: func(build api.PluginBuild) {
				build.OnLoad(api.OnLoadOptions{
					Filter: ".*",
				}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					err := watcher.Add(args.Path)
					return api.OnLoadResult{}, err
				})
			},
		}
		plugins = append(plugins, watchPlugin)
	}

	esbuildOpts := opts.Esbuild
	esbuildOpts.Plugins = plugins
	esbuildOpts.Incremental = opts.Watch

	if opts.Watch {
		for _, entrypoint := range esbuildOpts.EntryPoints {
			if !filepath.IsAbs(entrypoint) {
				entrypoint = filepath.Join(repo.RootDir, entrypoint)
			}
			if err := watcher.Add(entrypoint); err != nil {
				return fmt.Errorf("watching %q: %w", entrypoint, err)
			}
		}
	}

	result := api.Build(esbuildOpts)

	if opts.Types && opts.Package.Index != "" {
		args := []string{
			"--silent",
			"--no-banner",
		}
		for _, external := range esbuildOpts.External {
			args = append(args, "--external-imports="+external)
		}
		args = append(args,
			"--out-file", path.Join(repo.OutDir, "dist", opts.Package.Name, "index.d.ts"),
			opts.Package.Index,
		)
		cmd := exec.Command(
			path.Join(repo.RootDir, "node_modules", ".bin", "dts-bundle-generator"),
			args...,
		)
		cmd.Stdout = os.Stderr // Intentional redirect.
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("bundling type declarations: %w", err)
		}
	}

	g := new(errgroup.Group)

	abort := make(chan struct{})
	restart := make(chan struct{}, 1)

	g.Go(func() error {
		if len(result.Errors) > 0 {
			if !opts.Watch {
				return fmt.Errorf("build error")
			}
		}

		waitForChange := false
		for {
			proc := opts.CreateProcess()
			done := make(chan error, 1)

			buildOK := len(result.Errors) == 0
			shouldStart := buildOK && !waitForChange
			if shouldStart {
				if err := proc.Start(); err != nil {
					if !opts.Watch {
						return err
					}
					fmt.Fprintf(os.Stderr, "could not start: %v\n", err)
					waitForChange = true
				} else {
					go func() {
						done <- proc.Wait()
					}()
				}
			}
			select {
			case <-abort:
				if err := proc.Kill(); err != nil {
					fmt.Fprintf(os.Stderr, "could not kill: %v\n", err)
				}
				return nil
			case <-restart:
			loop:
				for {
					// Absorb extra restarts for a little while in case many files are changing at once.
					delay := time.After(50 * time.Millisecond)
					select {
					case <-restart:
					case <-delay:
						break loop
					}
				}
				if err := proc.Kill(); err != nil {
					fmt.Fprintf(os.Stderr, "could not kill: %v\n", err)
				}
				result = result.Rebuild()
				waitForChange = false
			case err := <-done:
				if !opts.Watch {
					return err
				}
				if err == nil {
					fmt.Fprintf(os.Stderr, "process finished\n")
				} else {
					fmt.Fprintf(os.Stderr, "process failure: %v\n", err)
				}
				waitForChange = true
			}
		}
	})

	if opts.Watch {
		g.Go(func() error {
			for {
				select {
				case _, ok := <-watcher.Events:
					if !ok {
						return nil
					}
					restart <- struct{}{}
				case err, ok := <-watcher.Errors:
					if !ok {
						close(abort)
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}
