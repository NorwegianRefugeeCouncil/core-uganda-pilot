import Constants from 'expo-constants';

export const linkingConfig = {
  prefixes: [
    `https://${Constants.manifest?.scheme}.com`,
    `${Constants.manifest?.scheme}://`,
  ],
  config: {
    screens: {
      Recipients: 'recipients',
      RecipientList: '/',
      RecipientRegistration: 'recipients/register',
      RecipientProfile: {
        path: '/:id',
        parse: {
          id: String,
        },
      },
    },
  },
};
