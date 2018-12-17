const express = require("express");
const next = require("next");
const compression = require("compression");
const cors = require("cors");
const helmet = require("helmet");
const proxy = require("http-proxy-middleware");
const port = parseInt(process.env.PORT, 10) || 3000;
const dev = process.env.NODE_ENV !== "production";
const app = next({ dev });
const handle = app.getRequestHandler();

global.fetch = require("isomorphic-fetch");
console.log('foo');

app.prepare().then(() => {
  const server = express();
  server.use(compression());
  server.use(cors());
  server.use(helmet());

  server.get("*", (req, res) => {
    return handle(req, res);
  });

  server.listen(port, err => {
    if (err) throw err;
    console.log(`
    From Skrop with ❤️
`);
  });
});
