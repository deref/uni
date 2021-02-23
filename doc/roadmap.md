# Roadmap

These are no particular order...

## Dependency Management

- hashing scheme instead of `--frozen` flag for `npm i` vs `npm ci`.

## Improved "engine checking"

Uni has support for _exact match_ of "engine" versions, but don't want exact
matches in published package.json `engines` field.

Need to support prefered production engine version, plus supported engine
ranges for published packages. These are two different things and need to
be treated differently.

## TypeScript

- Generate `tsconfig.json` file.
- `.d.ts` bundling. Use `tsc --watch` with `emitDeclarationOnly` for now.

## Package Organization

- common code split chunks -> common packages.
- support peer dependencies in built packages.

## Linting

- generate prettier config.
- generate eslint config.

## Execution

- crash on unhandled rejections for `run` command.
- on sigquit, dump open file handles, waiting promises, etc (is this possible w/o hurting runtime perf?)

## Import Paths

- provide something similar to Go's `internal` packages.
- disallow cyclic imports.
- automatically support absolute imports from `~`.
- disallow `..` in imports.
