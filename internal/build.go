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
		Name: "unirepo:deps",
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
		Outdir:    packageDir,
		Bundle:    true,
		Platform:  api.PlatformNode,
		Format:    api.FormatCommonJS,
		Write:     true,
		LogLevel:  api.LogLevelWarning,
		Sourcemap: api.SourceMapLinked,
		Plugins:   plugins,
		External:  getExternals(repo),
		Loader:    loaders,
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

	result := api.Build(buildOpts)

	pkgMetadata := PackageMetadata{
		Name:         pkg.Name,
		Private:      !pkg.Public,
		Description:  pkg.Description,
		Version:      opts.Version,
		Dependencies: dependencies,
		Bin:          bin,
	}

	if pkg.Index != "" {
		base := path.Base(pkg.Index)
		pkgMetadata.Main = strings.TrimSuffix(base, path.Ext(pkg.Index)) + ".js"
	}

	if err := WritePackageJSON(pkgMetadata, packageDir); err != nil {
		return err
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("build error")
	}

	return nil
}
