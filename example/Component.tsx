// XXX React/Browser not yet supported. This file is aspirational.

import React from "react";
import css from "./style.scss";

interface Props {}

export default function MyComponent(props: Props) {
  return <div className={css.red}>I am red.</div>;
}

export const main = () => {
  console.log("pointless main", <MyComponent />, css.main);
};
