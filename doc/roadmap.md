# Roadmap

- Generate `tsconfig.json` file.
- `.d.ts` bundling. Use `tsc --watch` with `emitDeclarationOnly` for now.
- common code split chunks -> common packages.
- disallow cyclic imports.
- automatically support absolute imports from `~`.
- disallow `..` in imports.
- provide something similar to Go's `internal` packages.
- generate prettier config.
- some way to determine which dependencies should be peer dependencies in generated packages.
- copy `engines` in to generated package.jsons
- hash the config file and if it hasn't changed, use `--immutable` deps install.
- crash on unhandled rejection for `run` command.
