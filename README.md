# monoclean

A clean mono workspace automation

# install

npm i monoclean

# Get stared

## declare you workspace `monoclean.yml`

```
engines:
  node: "v14.17.3"
  npm: "6.14.13"

workspace:
  name: "example"
  version: "1.0.0"

packages:
  "@example/a":
    folder: "a"
    dependencies:
      "@example/b": "*"

  "@example/b":
    folder: "b"
```

### `monoclean deps`

It will create and maintain automatically package.json, tsconfig.json, eslint.json, jest.config.js...

### `monoclean run`

`monoclean run [PACKAGE]` use esbuild to fasterly run `src/index.ts` on desired `[PACKAGE]`

### `monoclean`

`monoclean test [PACKAGE]` use esbuild and jest to fasterly test all packages in the workspace

### another commands
- [ ] monoclean build
- [ ] monoclean deploy

## Notes
- This project is in alfa version, use responsibly.
- This project is based on https://github.com/deref/uni
