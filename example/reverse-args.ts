import { reverse } from "~/util";

export const main = async (...args: string[]) => {
  for (const arg of reverse(args)) {
    console.log(arg);
  }
};
