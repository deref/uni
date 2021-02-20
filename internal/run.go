package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
)

type RunOptions struct {
	Watch      bool
	Entrypoint string
	Args       []string
}

// TODO: Need to handle interrupts in order to have a higher chance
// of cleaning up temporary files.

func Run(repo *Repository, opts RunOptions) error {
	if err := EnsureTmp(repo); err != nil {
		return err
	}

	dir, err := TempDir(repo, "run")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	script := `
		require('source-map-support').install();
		const { main } = require('./bundle.js');
		const args = process.argv.slice(2);
		void main(...args);
	`
	scriptPath := path.Join(dir, "script.js")
	if err := ioutil.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return err
	}

	var plugins []api.Plugin

	var watcher *fsnotify.Watcher
	if opts.Watch {
		var err error
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		watchPlugin := api.Plugin{
			Name: "unirepo:watch",
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

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{opts.Entrypoint},
		Outfile:     path.Join(dir, "bundle.js"),
		Bundle:      true,
		Platform:    api.PlatformNode,
		Format:      api.FormatCommonJS,
		Write:       true,
		LogLevel:    api.LogLevelWarning,
		Sourcemap:   api.SourceMapLinked,
		Incremental: opts.Watch,
		Plugins:     plugins,
	})

	g := new(errgroup.Group)

	abort := make(chan struct{})
	restart := make(chan struct{}, 1)

	g.Go(func() error {
		if len(result.Errors) > 0 {
			if !opts.Watch {
				return fmt.Errorf("build error")
			}
		}

		for {
			nodeArgs := append([]string{scriptPath}, opts.Args...)
			node := exec.Command("node", nodeArgs...)
			node.Stdin = os.Stdin
			node.Stdout = os.Stdout
			node.Stderr = os.Stderr
			done := make(chan error, 1)
			if err := node.Start(); err != nil {
				if !opts.Watch {
					return err
				}
				fmt.Fprintf(os.Stderr, "could not start: %v", err)
			} else {
				go func() {
					done <- node.Wait()
				}()
			}
			select {
			case <-abort:
				if err := node.Process.Kill(); err != nil {
					fmt.Fprintf(os.Stderr, "could not kill: %v", err)
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
				if err := node.Process.Kill(); err != nil {
					fmt.Fprintf(os.Stderr, "could not kill: %v", err)
				}
				result = result.Rebuild()
			case err := <-done:
				if !opts.Watch {
					return err
				}
				fmt.Fprintf(os.Stderr, "process terminated: %v", err)
			}
		}
	})

	if opts.Watch {
		g.Go(func() error {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return nil
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						restart <- struct{}{}
					}
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
