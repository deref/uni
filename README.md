# Unirepo

Unirepo is an extremely opinionated TypeScript build tool.

Typical monorepo management tools in the Node.js ecosystem provide automation
around package maintenance, but still permit and require users to muck around
with poly-package configuration. Package boundaries must be manually
maintained, each with its own sub-configuration.

Unirepo is different because sub-package configuration is managed centrally and
uniformly. Package boundaries are managed automatically via bundling and code
splitting.

You will need one and only one package.json file. Same for tsconfig.json.

Unirepo is _fast_ because it is ships as a native binary and builds your code
using [esbuild][1].

Additionally, Unirepo has a `run` subcommand that acts as a substitute for both
[`ts-node`][2]. The `run` subcomain also supports a `--watch` flag, and so acts
as a substitute for [`node-dev`][3] (or [`ts-node-dev`][4]) as well.

## Installation

TODO: Published binaries. Until then, install the package with `go get` or similar.

## Configuration

See the [example](./example).

## Status

Alpha! Don't use this yet.

- TODO: .d.ts bundling.
- TODO: common code split chunks -> common packages.

[1]: https://esbuild.github.io/
[2]: https://github.com/TypeStrong/ts-node
[3]: https://github.com/fgnass/node-dev
[4]: https://github.com/wclr/ts-node-dev
