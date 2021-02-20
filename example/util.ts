import { formatISO } from "date-fns";

export const first = <T>(xs: T[]): T => {
  if (xs.length === 0) {
    throw Error("first of empty array");
  }
  return xs[0];
};

export const last = <T>(xs: T[]): T => {
  const n = xs.length;
  if (n === 0) {
    throw Error("last of empty array");
  }
  return xs[n - 1];
};

export const reverse = <T>(xs: T[]): T[] => {
  const ys = xs.slice();
  return ys.reverse();
};

export const nowISO = () => formatISO(new Date());
