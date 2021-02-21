# Versioning

This project uses `$SOME_CUTE_VERSIONING_SCHEME_NAME`. It works like this:

Version numbers have the short form:

`<major>.<minor>`

## Stability

Functionality can be either **stable** or **unstable**. Stable features are not
subject to breaking changes, even in major version upgrades.

Unstable features fall in to four categories: **alpha**, **beta**,
**deprecated**, and **internal**.

**Alpha** features are subject to breaking changes at anytime and without
warning.

**Beta** features _may_ experience breaking changes, but the goal is for them to
stabilize, and so such breakage will be accompanied by guidance.

**Deprecated** features will be _removed_ in the next major version release.
Warning and a migration plan will be provided.

**Internal** features are also subject to breaking changes at anytime. They
differ from _alpha_ features in that they are not expected to stabilize.

## Version Numbering

When the major version is 0, it is called the "alpha release". All features are
considered _alpha_, unless explicitly declared otherwise. Even if designated
_stable_, functionality in an alpha release should should be treated as _beta_
at best.

For major version 1 and up, all _documented_ features are considered _stable_,
unless otherwise indicated. All _undocumented_ functionality is considered
_internal_.
