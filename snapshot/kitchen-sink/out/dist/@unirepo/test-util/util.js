var __defProp = Object.defineProperty;
var __markAsModule = (target) => __defProp(target, "__esModule", {value: true});
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, {get: all[name], enumerable: true});
};

// util.ts
__markAsModule(exports);
__export(exports, {
  first: () => first,
  last: () => last
});
var first = (xs) => {
  if (xs.length === 0) {
    throw Error("first of empty array");
  }
  return xs[0];
};
var last = (xs) => {
  const n = xs.length;
  if (n === 0) {
    throw Error("last of empty array");
  }
  return xs[n - 1];
};
//# sourceMappingURL=util.js.map
