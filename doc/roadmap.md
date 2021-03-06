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

Also needed: _self_ engine check. Validating current version of Unirepo.

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

## REPL

`uni repl` should support TypeScript syntax.

## Bin Stubs

When specifying dependencies, it should be possible to specify a subset
of executables for which stub shell scripts would be generated and placed
in to `./bin`. These executables would check engines before executing.

It might also be desirable to create `node` and `npm` binstubs to alleviate the
need for a node version manager.
