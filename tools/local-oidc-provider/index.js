const https = require("https")
const http = require("http")
const fs = require("fs");
const express = require('express');
const Provider = require('oidc-provider');

let PORT = 3000
if (process.env.PORT) {
    PORT = process.env.PORT
}

let issuer = `http://localhost:${PORT}`
if (process.env.ISSUER) {
    issuer = process.env.ISSUER
}

let config = {}
if (process.env.CONFIG_FILE) {
    config = require(process.env.CONFIG_FILE)
}

config.pkce = {}
config.pkce.required = () => false

config.renderError = async (ctx, out, error) => {
    ctx.type = 'html';
    ctx.body = `<!DOCTYPE html>
    <head>
      <title>oops! something went wrong</title>
      <style>/* css and html classes omitted for brevity, see lib/helpers/defaults.js */</style>
    </head>
    <body>
      <div>
        <h1>oops! something went wrong</h1>
        ${Object.entries(out).map(([key, value]) => `<pre><strong>${key}</strong>: ${value}</pre>`).join('')}
        <hr/>
        ${error}
      </div>
    </body>
    </html>`;
}
const app = express();

if (process.env.TLS_KEY || process.env.TLS_CERT){
    https.createServer({
        key: fs.readFileSync(process.env.TLS_KEY),
        cert: fs.readFileSync(process.env.TLS_CERT),
    }, app).listen(PORT)
} else {
    http.createServer(app).listen(PORT)
}

console.log("config", config)

const provider = new Provider(issuer, config);
app.use(provider.callback())
