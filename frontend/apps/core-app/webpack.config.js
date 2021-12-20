// eslint-disable-next-line @typescript-eslint/no-var-requires
const { createWebpackConfigAsync } = require('expo-yarn-workspaces/webpack');

module.exports = async function (env, argv) {
  const config = await createWebpackConfigAsync(env, argv);
  return config;
};
