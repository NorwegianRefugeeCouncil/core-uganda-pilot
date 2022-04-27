import Constants from 'expo-constants';
import * as Linking from 'expo-linking';

export const linkingConfig = {
  prefixes: [
    Linking.createURL('/'),
    `https://${Constants.manifest?.scheme}.com`,
  ],
};
