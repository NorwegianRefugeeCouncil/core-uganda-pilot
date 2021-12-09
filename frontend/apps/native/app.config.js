import 'dotenv/config';

export default ({ config }) => {
  switch (process.env.NODE_ENV) {
    case 'development':
    case 'dev':
      return {
        ...config,
        plugins: ['@config-plugins/android-jsc-intl'],
        extra: {
          server_uri: process.env.REACT_APP_SERVER_URL,
          client_id: process.env.REACT_APP_OAUTH_CLIENT_ID,
          issuer: process.env.REACT_APP_OIDC_ISSUER,
          scopes: process.env.REACT_APP_OAUTH_SCOPE.split(' '),
        },
      };
    case 'production':
    case 'prod':
      return config;
    case 'staging':
      return config;
    default:
      return config;
  }
};
