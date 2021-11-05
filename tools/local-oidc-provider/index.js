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

config.renderError = async(ctx, out, error) => {
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

console.log("config", config)

const oidc = new Provider(issuer, config);
oidc.listen(PORT, () => {
    console.log(`issuer ${issuer}`)
    console.log(`port ${PORT}`)
})
