// Cases to handle
//
// without watch
//   build error
//   build ok
//     program fails to start
//     program terminates success
//     program terminates failure
//
// with watch
//   build error
//   build ok                              waitForChange
//     program fails to start                  true      wait for change
//     program terminates prematurely          true      wait for change
//     code changes                            false     restart
//     interrupt                               false     exi

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
	BuildOnly  bool
}

// TODO: Need to handle interrupts in order to have a higher chance
// of cleaning up temporary files.

// Status code may be returend within an exec.ExitError return value.
func Run(repo *Repository, opts RunOptions) error {
	if err := EnsureTmp(repo); err != nil {
		return err
	}

	dir, err := TempDir(repo, "run")
	if err != nil {
		return err
	}
	if !opts.BuildOnly {
		defer os.RemoveAll(dir)
	}

	// See also `shim` in Build.
	script := fmt.Sprintf(`require('source-map-support').install();
const { main } = require('./bundle.js');
if (typeof main === 'function') {
	const args = process.argv.slice(2);
	void (async () => {
		const exitCode = await main(...args);
		process.exit(exitCode ?? 0);
	})();
} else {
	process.stdout.write('error: %s does not export a main function\n', () => {
		process.exit(1);
	});
}
`, opts.Entrypoint)
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
		AbsWorkingDir: repo.RootDir,
		EntryPoints:   []string{opts.Entrypoint},
		Outfile:       path.Join(dir, "bundle.js"),
		Bundle:        true,
		Platform:      api.PlatformNode,
		Format:        api.FormatCommonJS,
		Write:         true,
		LogLevel:      api.LogLevelWarning,
		Sourcemap:     api.SourceMapLinked,
		Incremental:   opts.Watch,
		Plugins:       plugins,
		External:      getExternals(repo),
		Loader:        loaders,
	})

	if opts.BuildOnly {
		fmt.Println(dir)
		return nil
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
			nodeArgs := append([]string{scriptPath}, opts.Args...)
			node := exec.Command("node", nodeArgs...)
			node.Stdin = os.Stdin
			node.Stdout = os.Stdout
			node.Stderr = os.Stderr
			done := make(chan error, 1)

			buildOK := len(result.Errors) == 0
			shouldStart := buildOK && !waitForChange
			if shouldStart {
				if err := node.Start(); err != nil {
					if !opts.Watch {
						return err
					}
					fmt.Fprintf(os.Stderr, "could not start: %v\n", err)
					waitForChange = true
				} else {
					go func() {
						done <- node.Wait()
					}()
				}
			}
			select {
			case <-abort:
				if err := node.Process.Kill(); err != nil {
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
				proc := node.Process
				if proc != nil {
					if err := proc.Kill(); err != nil {
						fmt.Fprintf(os.Stderr, "could not kill: %v\n", err)
					}
				}
				result = result.Rebuild()
				waitForChange = false
			case err := <-done:
				if !opts.Watch {
					return err
				}
				fmt.Fprintf(os.Stderr, "process terminated: %v\n", err)
				waitForChange = true
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
