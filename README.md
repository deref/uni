# Unirepo

Unirepo is an extremely opinionated TypeScript build tool.

Typical monorepo management tools in the Node.js ecosystem provide automation
around package maintenance, but still permit and require users to muck around
with poly-package configuration. Package boundaries must be manually
maintained, each with its own sub-configuration.

Unirepo is different because sub-package configuration is managed centrally and
uniformly. Package boundaries are managed automatically via bundling and code
splitting.

You will have one and only one `package.json` file. Same for `tsconfig.json`.
Your `package.json` file will be generated from a source configuration file
that - believe it or not - allows code comments.

Unirepo is _fast_ because it is ships as a native binary and builds your code
using [esbuild][1].

Additionally, Unirepo has a `run` subcommand that acts as a substitute for
[`ts-node`][2]. The `run` subcomain also supports a `--watch` flag, and so acts
as a substitute for [`node-dev`][3] (or [`ts-node-dev`][4]) as well. Sourcemaps
are always enabled.

Want to see it in action?
Check out the [Demo Video](https://www.youtube.com/watch?v=RJfLA7EM-Uw)!

## Status

Alpha! Don't use this yet.

- TODO: `.d.ts` bundling. Use `tsc --watch` with `emitDeclarationOnly` for now.
- TODO: common code split chunks -> common packages.
- TODO: disallow cyclic imports.
- TODO: support absolute imports from `~`.
- TODO: disallow `..` in imports.
- TODO: provide something similar to Go's `internal` packages.

## Installation

TODO: Published binaries. Until then, install the package with `go get` or similar.

## Configuration

See the [example](./example).

## Usage

### Setup

1. Create a `uni.yml` file with some package entrypoints.
2. Manually add dependencies to your config file.
3. Run `uni deps`.

### Development

- Use `uni run src/program.ts` to execute programs. They must export a `main` function.
- Use `uni build some-package` to pre-compile into `out/dist`.

### Publishing

Here's the steps to do in your CI flow:

1. `uni build --version $VERSION` to create packages with version numbers.
2. `uni pack` to create packed `.tgz` files.
3. `uni publish` to automate `npm publish ./path/to/package.tgz`.

[1]: https://esbuild.github.io/
[2]: https://github.com/TypeStrong/ts-node
[3]: https://github.com/fgnass/node-dev
[4]: https://github.com/wclr/ts-node-dev
