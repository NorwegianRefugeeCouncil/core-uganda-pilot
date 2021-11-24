const { build } = require('esbuild');
const { exec } = require('child_process')
const path = require('path');
const { pnpPlugin } = require('@yarnpkg/esbuild-plugin-pnp');

async function buildTs() {
    return await new Promise((resolve, reject) => {
        exec('yarn tsc -b -f', (error, stdout) => {
            console.log(stdout);
            if (error) {
                reject('failed to build ts projects');
                return
            }
            resolve('done!')
        })
    })
}

const baseLibConfig = {
    bundle: true,
    format: 'esm',
    target: ['es6'],
    plugins: [pnpPlugin()],
};

async function buildAuthClient(options = { env: 'development' }) {
    const { env } = options;

    const p = path.resolve('libs/auth-client');
    console.log('>> Building auth-client')
    try {
        await build({
            ...baseLibConfig,
            entryPoints: [path.resolve(p, 'src/index.ts')],
            outdir: path.join(p, 'dist'),
            define: {
                'process.env.NODE_ENV': `"${env}"`
            },
            minify: env === 'production',
            sourcemap: env === 'development'
        });
    } catch (e) { throw e }
}

async function buildApiClient(options = { env: 'development' }) {
    const { env } = options;

    const p = path.resolve('libs/api-client');
    console.log('>> Building api-client')
    try {
        await build({
            ...baseLibConfig,
            entryPoints: [path.resolve(p, 'src/index.ts')],
            outdir: path.join(p, 'dist'),
            define: {
                'process.env.NODE_ENV': `"${env}"`
            },
            minify: env === 'production',
            sourcemap: env === 'development'
        });
    } catch (e) { throw e }
}

async function buildLibs() {
    await Promise.all([
        buildAuthClient(),
        buildApiClient(),
    ]).catch(err => { throw err })
}

(async function() {
    try {
        console.log('>> Compiling typescript (declarations)')
        await buildTs();
        console.log('>> Building libraries')
        await buildLibs();
    } catch (e) {
        console.error(e)
        process.exit()
    } finally {
        console.log('>> Librairies built successfully! 🎉🎉🎉')
    }
})()
