package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
)

type BuildOptions struct {
	Package *Package
	Version string
	Types   bool
	Watch   bool
}

func Build(repo *Repository, opts BuildOptions) error {
	pkg := opts.Package

	packageDir := path.Join(repo.OutDir, "dist", pkg.Name)
	if err := os.RemoveAll(packageDir); err != nil {
		return err
	}
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		return err
	}

	if err := EnsureTmp(repo); err != nil {
		return err
	}

	var mx sync.Mutex
	dependencies := make(map[string]string)

	depPrefix := path.Join(repo.RootDir, "node_modules")
	isFileFromDeps := func(filepath string) bool {
		return strings.HasPrefix(filepath, depPrefix)
	}

	depsPlugin := api.Plugin{
		Name: "monoclean:deps",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{
				Filter: ".*",
			}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				if isFileFromDeps(args.Importer) {
					return api.OnResolveResult{}, nil
				}
				moduleName := args.Path
				if dependency, ok := repo.Dependencies[moduleName]; ok {
					mx.Lock()
					dependencies[moduleName] = dependency.Version
					mx.Unlock()
				}
				return api.OnResolveResult{}, nil
			})
		},
	}

	plugins := []api.Plugin{
		depsPlugin,
	}

	indexPath := path.Join(repo.RootDir, pkg.Index)

	buildOpts := api.BuildOptions{
		AbsWorkingDir: repo.RootDir,
		Outdir:        packageDir,
		Bundle:        true,
		Platform:      api.PlatformNode,
		Format:        api.FormatCommonJS,
		Write:         true,
		LogLevel:      api.LogLevelWarning,
		Sourcemap:     api.SourceMapLinked,
		Plugins:       plugins,
		External:      getExternals(repo),
		Loader:        loaders,
		// TODO: Splitting: true,
	}

	if pkg.Index != "" {
		buildOpts.EntryPoints = append(buildOpts.EntryPoints, indexPath)
	}

	bin := make(map[string]string)
	for executableName, executable := range pkg.Executables {
		buildOpts.EntryPoints = append(buildOpts.EntryPoints, executable.Entrypoint)
		bin[executableName] = executableName

		entrypointExt := path.Ext(executable.Entrypoint)
		entrypointOut := strings.TrimSuffix(path.Base(executable.Entrypoint), entrypointExt) + ".js"

		// See also `script` in Run.
		shim := fmt.Sprintf(`#!/usr/bin/env node

const { inspect } = require('util');
process.on('uncaughtException', (exception) => {
  process.stderr.write('uncaught exception: ' + inspect(exception) + '\n', () => {
    process.exit(1);
  });
});
process.on('unhandledRejection', (reason, promise) => {
  process.stderr.write(
    'unhandled rejection at: ' + inspect(promise) + '\nreason: ' + inspect(reason) + '\n',
    () => {
      process.exit(1);
    },
  );
})

const { main } = require('./%s');
const args = process.argv.slice(2);
void (async () => {
	const exitCode = await main(...args);
	process.exit(exitCode ?? 0);
})();
`, entrypointOut)

		executablePath := path.Join(packageDir, executableName)
		if err := ioutil.WriteFile(executablePath, []byte(shim), 0755); err != nil {
			return err
		}
	}

	isScoped := strings.HasPrefix(pkg.Name, "@")
	private := !(pkg.Public || isScoped)

	return buildAndWatch{
		Repository: repo,
		Esbuild:    buildOpts,
		Types:      opts.Types,
		Watch:      opts.Watch,
		Package:    pkg,
		CreateProcess: func() process {
			return &funcProcess{
				start: func() error {
					pkgMetadata := PackageMetadata{
						Name:         pkg.Name,
						Private:      private,
						Description:  pkg.Description,
						Version:      opts.Version,
						Dependencies: dependencies,
						Bin:          bin,
						Repository:   repo.Url,
						PublishConfig: &PublishConfig{
							Registry: repo.Registry,
						},
					}

					if pkg.Index != "" {
						base := path.Base(pkg.Index)
						pkgMetadata.Main = strings.TrimSuffix(base, path.Ext(pkg.Index)) + ".js"
					}

					if err := WritePackageJSON(pkgMetadata, packageDir); err != nil {
						return err
					}

					return nil
				},
			}
		},
	}.Run()
}

type funcProcess struct {
	start func() error
}

func (proc *funcProcess) Start() error {
	return proc.start()
}

func (proc *funcProcess) Kill() error {
	return nil
}

func (proc *funcProcess) Wait() error {
	return nil
}
