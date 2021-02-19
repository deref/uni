package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func Build() error {
	decl := PackageDeclaration{
		Name:        "@unirepo/example-util",
		Description: "An example package of utility functions.",
		Entrypoint:  "src/util.ts",
	}
	packageDir := path.Join(outDir, "node_modules", decl.Name)

	if err := os.MkdirAll(packageDir, 0755); err != nil {
		return err
	}

	if err := EnsureTmp(); err != nil {
		return err
	}

	// XXX get from configs.
	allDependencies := map[string]string{
		"date-fns": "2.17.0",
	}

	dependencies := make(map[string]string)

	depPrefix := "/Users/brandonbloom/Projects/unirepo/example/node_modules/"
	isFileFromDeps := func(filepath string) bool {
		return strings.HasPrefix(filepath, depPrefix)
	}

	var buildPlugin = api.Plugin{
		Name: "unirepo",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{
					Filter: `.*`,
				},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					if isFileFromDeps(args.Importer) {
						return api.OnResolveResult{}, nil
					}
					moduleName := args.Path
					if version, ok := allDependencies[moduleName]; ok {
						dependencies[moduleName] = version
					}
					return api.OnResolveResult{}, nil
				},
			)
		},
	}

	mainRelpath := "index.cjs.js"
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{decl.Entrypoint},
		Outfile:     path.Join(packageDir, mainRelpath),
		Bundle:      true,
		Platform:    api.PlatformNode,
		Format:      api.FormatCommonJS,
		Write:       true,
		Plugins: []api.Plugin{
			buildPlugin,
		},
	})

	pkgMetadata := PackageMetadata{
		Name:         decl.Name,
		Description:  decl.Description,
		Main:         mainRelpath,
		Dependencies: dependencies,
	}
	if err := WritePackageJSON(pkgMetadata, packageDir); err != nil {
		return err
	}

	if len(result.Errors) > 0 {
		// XXX better reporting.
		for _, err := range result.Errors {
			fmt.Println(err)
		}
		return fmt.Errorf("build error")
	}

	return nil
}
