# Anti-features

A list of things you cannot do with Unirepo, why, and what you should do
instead.

## Mixed-version Dependencies

In Unirepo, you may only explicitly depend on a single version of any given
package.

Note that this doesn't mean that you can't have multiple versions of a
dependency in your package tree. Transitive dependency resolution is delegated
to the underlying package manager, and so is still subject to that complexity.

If for some reason you need to directly depend on two incompatible versions of
a dependency, treat that as two different dependencies! Fork the project,
publish it under a new name, and import that new name. This policy is inspired
by [Semantic Import Versioning][1] as pioneered by the Go community.

[1]: https://research.swtch.com/vgo-import

## Development Dependencies

The distinction between production and development dependencies is maintained
automatically and implicitly. If you publish an entrypoint that directly or
indirectly depends on another package, that package is considered a production
dependency of the published package. Any dependency that is not required by a
published package can be logically considered a development dependency.

## Automated Dependency Changes

While undoubtedly convenient for common tasks, commands that edit your
dependency list automatically are only generally possible in environments
where the dependency list is treated as a data file instead of source code.
Since the Unirepo configuration file may contain comments and other structure,
a command similar to `npm install --save` would potentially corrupt whatever
organizational scheme has been imposed. Instead, use `npm info` to find the
latest version number and add it to the dependency list manually. And maybe
include a comment too!

## Import Cycles

Module import cycles are explicitly disallowed. Cyclic dependencies lead to a
number of unnecessary complications. Mutual recursion is frequently essential
to a design, but import cycles are accidental complexity. Any fundamental knots
can be untangled by hoisting them to a central location and introducing
indirection. This indirection typically takes the form of interface and
dependency injection, or mutation during initialization.
