export const main = async (...args: string[]) => {
  for (const arg of args.reverse()) {
    console.log(arg);
  }
};
