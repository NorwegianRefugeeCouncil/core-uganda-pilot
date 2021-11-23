const { build } = require('esbuild');
const { spawn } = require('child_process')
const path = require('path');

export async function buildTs() {
    return await new Promise((resolve, reject) => {
        const tsc = spawn('tsc', ['-b', '-f'], { shell: true })
        tsc.stdout.on('data', (data: Buffer) => console.log(data.toString()))
        tsc.stderr.on('data', (data: Buffer) => console.error(data.toString()))
        tsc.on('close', code => {
            if (code === 0) resolve('done');
            else reject(new Error('errors encountered trying while building ts projects'));
        })
    })
}

type BuildOptions = {
    env: 'production' | 'development'
};

const baseLibConfig = {
    entryPoints: ['src/index.ts'],
    outdir: 'dist',
    tsconfig: 'tsconfig.json',
    bundle: true,
    target: ['es6', 'node16'],
};

export async function buildAuthClient(options: BuildOptions = { env: 'development' }) {
    const { env } = options;

    console.log('>> Building auth-client')
    try {
        await build({
            ...baseLibConfig,
            absWorkingDir: path.join(process.cwd(), 'libs/auth-client'),
            define: {
                'process.env.NODE_ENV': `"${env}"`
            },
            minify: env === 'production',
            sourcemap: env === 'development'
        });
    } catch (e) { throw e }
}

export async function buildApiClient(options: BuildOptions = { env: 'development' }) {
    const { env } = options;

    console.log('>> Building api-client')
    try {
        await build({
            ...baseLibConfig,
            absWorkingDir: path.join(process.cwd(), 'libs/api-client'),
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
