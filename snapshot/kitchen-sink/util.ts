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
