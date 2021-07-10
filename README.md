# Unirepo

Unirepo is an extremely opinionated TypeScript build tool.

Typical monorepo management tools in the Node.js ecosystem provide automation
around package maintenance, but still permit and require users to muck around
with poly-package configuration. Package boundaries must be manually
maintained, each with its own sub-configuration.

Unirepo is different because sub-package configuration is managed centrally and
uniformly. Package boundaries are managed automatically via bundling and code
splitting.

You will have one and only one list of dependencies. Your `package.json` files
will be generated from that source configuration. Same for `tsconfig.json`.
The source configuration file - believe it or not - allows code comments.

Unirepo is _fast_ because it is ships as a native binary and builds your code
using [esbuild][1].

Additionally, Unirepo has a `run` subcommand that acts as a substitute for
[`ts-node`][2]. The `run` subcommand also supports a `--watch` flag, and so
acts as a substitute for [`node-dev`][3] (or [`ts-node-dev`][4]) as well.
Sourcemaps are always enabled.

As mentioned, Unirepo is extremely opinionated. Those opinions will evolve into
documentation, including a growing list of
[anti-features](./doc/anti-features.md).

Want to see it in action?
Check out the [Demo Video](https://www.youtube.com/watch?v=RJfLA7EM-Uw)!

## Status

Alpha! Don't use this yet.

See the [versioning guide](./doc/versioning.md) and the
[roadmap](./doc/roadmap.md).

Only works for targeting Node currently. Targeting Browsers is planned.

## Installation

```bash
go get -u github.com/deref/uni
```

## Usage

- [Configuration](./doc/config.md)
- [Migration Guide](./doc/migrate.md)

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

## Other Features

### Patching

The [patch-package][5] utility is always available.

### Engine Checking

Functionality similar to [check-engine][6] is builtin, but much faster
and with caching.

### Executables

Any runnable script can be exposed as an executable in a package. A shim script
(with a `#!`) will be produced automatically.

[1]: https://esbuild.github.io/
[2]: https://github.com/TypeStrong/ts-node
[3]: https://github.com/fgnass/node-dev
[4]: https://github.com/wclr/ts-node-dev
[5]: https://github.com/ds300/patch-package
[6]: https://github.com/mohlsen/check-engine
