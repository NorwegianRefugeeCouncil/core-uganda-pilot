import { build } from 'esbuild';

type BuildOptions = {
    env: 'production' | 'development'
};

const baseConfig = {
    bundle: true,
    target: ['es6', 'node16'],
};

export async function buildAuthClient(options: BuildOptions = { env: 'development' }) {
    const { env } = options;

    await build({
        ...baseConfig,
        entryPoints: ['libs/auth-client/src/index.ts'],
        outdir: 'libs/auth-client/dist',
        define: {
            'process.env.NODE_ENV': `"${env}"`
        },
        minify: env === 'production',
        sourcemap: env === 'development'
    });
}

export async function buildApiClient(options: BuildOptions = { env: 'development' }) {
    const { env } = options;

    await build({
        ...baseConfig,
        entryPoints: ['libs/api-client/src/index.ts'],
        outdir: 'libs/api-client/dist',
        define: {
            'process.env.NODE_ENV': `"${env}"`
        },
        minify: env === 'production',
        sourcemap: env === 'development'
    });
}

export async function buildDesignSystem(options: BuildOptions = { env: 'development' }) {
    const { env } = options;

    await build({
        ...baseConfig,
        entryPoints: ['libs/design-system/src/index.ts'],
        outdir: 'libs/design-system/dist',
        define: {
            'process.env.NODE_ENV': `"${env}"`
        },
        minify: env === 'production',
        sourcemap: env === 'development'
    });
}

async function buildAll() {
    await Promise.all([
        buildAuthClient(),
        buildApiClient(),
        buildDesignSystem()
    ]);
}

buildAll();
