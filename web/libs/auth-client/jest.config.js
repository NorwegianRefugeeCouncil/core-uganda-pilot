/** @type {import('ts-jest/dist/types').InitialOptionsTsJest} */
module.exports = {
    globals: {
        'ts-jest': {
            // ts-jest configuration goes here
    preset: 'ts-jest/presets/js-with-babel-esm',
    testEnvironment: 'node',
            useESM: true
        },
    },
};
