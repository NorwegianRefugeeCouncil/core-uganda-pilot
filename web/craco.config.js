const path = require('path');
const { getLoader, loaderByName } = require('@craco/craco');

const packages = [];
packages.push(path.resolve(__dirname, 'src/'));

module.exports = {
  webpack: {
    configure: (webpackConfig) => {
      const { isFound, match } = getLoader(webpackConfig, loaderByName('babel-loader'));
      if (isFound) {
        const include = Array.isArray(match.loader.include) ? match.loader.include : [match.loader.include];

        match.loader.include = include.concat(packages);
      } else {
          throw new Error("failed to find babel-loader");
      }
      return webpackConfig;
    },
  },
};
