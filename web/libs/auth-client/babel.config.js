module.exports = function (api) {
    const presets = [
        ['@babel/preset-env', {targets: {node: 'current'}}],
        // '@babel/preset-typescript',
    ];

    // const presets = ['@babel/preset-env', '@babel/preset-react'];

    const plugins = [
        // '@babel/plugin-transform-runtime',
        '@babel/plugin-syntax-jsx',
    ];
    //
    /** this is just for minimal working purposes,
     * for testing larger applications it is
     * advisable to cache the transpiled modules in
     * node_modules/.bin/.cache/@babel/register* */
    api.cache(false);

    return {
        presets,
        plugins
    };
};
