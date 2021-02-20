import http from "http";

export const main = async (...args: string[]) => {
  const port = 3001;
  http
    .createServer((req, res) => {
      res.write("Hello there!\n");
      res.end();
    })
    .listen(port);
  console.log("listening on port:", port);
};
