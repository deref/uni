import { color as theirs } from "cycle-green";

export const color = "blue";

export const main = async () => {
  console.log("ours:", color);
  console.log("theirs:", theirs);
  throw Error("this should be unreachable with uni run");
};
