# Config

Configuration is specified in a `uni.yml` file at the project root.

For an example project, see the [`uni.yml` file](../example/uni.yml) in the
[example directory](../example).

# `repository`

URL of the containing code repository. If provided, this property is copied
into all generated `package.json` files.

# `registry`

Where to publish packages to. Defaults to `https://registry.npmjs.org/`.

# `packages`

Map of packages to be published, keyed by name.

## `packages.<package-name>`

`package-name` is of the form `name` or `@scope/name` as specified
by NPM.

### `packages.<package-name>.index`

Path to the code file that exports the public interface of the package.

### `packages.<package-name>.executables.<executable-name>: <entrypoint>`

Map of executables to be included in the package.

`executable-name` is the filename to use for an executable program to be
included in the built package.

`entrypoint` is the path to an entrypoint module which is suitable for use
with `uni run` (ie. it must export a `main` function).

### `packages.<package-name>.public`

_Default:_ `false`

Setting to true will allow packages to be published to a public registry for
anyone to download.

### `packages.<package-name>.description`

A short description to accompany the package name when published to a registry.

# `engines`

Specifies required external programs versions. If provided, these are checked before running any operations that are sensitive to these programs.

NOTE: These engines are not copied to any generated package.json files yet.

## `engines.<engine-name>: <version>`

`engine-name` must be from the following known list of supported engines:

- `node`
- `npm`

`version` is the exactly version number expected.

# `dependencies`

A map of dependency package names to their version numbers.

## `dependencies.<dependency-name>: <version>`

This is the package name and version number to depend on, as specified by NPM.

The version number will passthrough to NPM unmodified, but this is an
implementation detail and may change. Therefore, you should avoid using version
ranges or specifiers like `^`.
