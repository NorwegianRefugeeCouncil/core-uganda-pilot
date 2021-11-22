const config = {
    transform: {
        "\\.[jt]sx?$": ['esbuild-jest', {
            loaders: {
                '.spec.js': 'jsx',
                '.test.js': 'jsx',
                '.js': 'jsx'
            }
        }
        ]
    },
    testPathIgnorePatterns: ['/dist/', '/types/'],
    moduleFileExtensions: ['ts', 'tsx', 'jsx', 'js', 'json', 'node']
}


module.exports = config;
