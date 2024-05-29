import { Platform } from 'react-native';
import * as AuthSession from 'expo-auth-session';
import { openAuthSessionAsync } from 'expo-web-browser';
import Constants from 'expo-constants';

import { formsClient } from '../../clients/formsClient';

const authorizationEndpoint = `${Constants.manifest?.extra?.issuer}/oauth2/sessions/logout`;
const redirectUri = AuthSession.makeRedirectUri({ useProxy: false });

export const logout = async () => {
  try {
    await openAuthSessionAsync(
      `${authorizationEndpoint}?post_logout_redirect=${redirectUri}`,
      'redirectUrl',
    );
  } catch (err) {
    console.error(err);
  } finally {
    if (Platform.OS === 'web') {
      sessionStorage.removeItem('tokenResponse');
    }
    formsClient.setToken('');
  }
};
