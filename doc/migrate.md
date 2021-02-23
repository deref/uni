# Migration Guide

Migrating from _%SOME_TOOL%_ to Unirepo?

You might need to make some changes.

## ts-node, node-dev, ts-node-dev

Instead of `ts-node ./path/to/script.ts`, use `uni run ./path/to/module.ts`.

The module must export a `main` function. For example:

```typescript
export const main = async (...args: string[]) => {
  console.log("see uni run");
  return 0; // Return an exit code (optional).
};
```

This is because `uni run` does not make a distinction between modules and
scripts; there are only modules. It should always be safe and side-effect free
to load a (trusted) module.
