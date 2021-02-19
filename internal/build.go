package internal

import (
	"fmt"
	"os"
	"path"

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

	mainRelpath := "index.cjs.js"
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{decl.Entrypoint},
		Outfile:     path.Join(packageDir, mainRelpath),
		Platform:    api.PlatformNode,
		Format:      api.FormatCommonJS,
		Write:       true,
	})

	dependencies := map[string]string{} // XXX fill me.

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
