var __create = Object.create;
var __defProp = Object.defineProperty;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __markAsModule = (target) => __defProp(target, "__esModule", {value: true});
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, {get: all[name], enumerable: true});
};
var __exportStar = (target, module2, desc) => {
  if (module2 && typeof module2 === "object" || typeof module2 === "function") {
    for (let key of __getOwnPropNames(module2))
      if (!__hasOwnProp.call(target, key) && key !== "default")
        __defProp(target, key, {get: () => module2[key], enumerable: !(desc = __getOwnPropDesc(module2, key)) || desc.enumerable});
  }
  return target;
};
var __toModule = (module2) => {
  if (module2 && module2.__esModule)
    return module2;
  return __exportStar(__markAsModule(__defProp(module2 != null ? __create(__getProtoOf(module2)) : {}, "default", {value: module2, enumerable: true})), module2);
};

// util.ts
__markAsModule(exports);
__export(exports, {
  first: () => first,
  last: () => last,
  nowISO: () => nowISO,
  reverse: () => reverse
});
var import_date_fns = __toModule(require("date-fns"));
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
var reverse = (xs) => {
  const ys = xs.slice();
  return ys.reverse();
};
var nowISO = () => import_date_fns.formatISO(new Date());
//# sourceMappingURL=util.js.map
