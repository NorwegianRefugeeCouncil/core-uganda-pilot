import { Platform } from 'react-native';
import * as AuthSession from 'expo-auth-session';
import { openAuthSessionAsync } from 'expo-web-browser';
import Constants from 'expo-constants';

const authorizationEndpoint = `${Constants.manifest?.extra?.issuer}/oauth2/sessions/logout`;
const redirectUri = AuthSession.makeRedirectUri({ useProxy: false });

export const logout = async () => {
  try {
    await openAuthSessionAsync(
      `${authorizationEndpoint}?post_logout_redirect=${redirectUri}`,
      'redirectUrl',
    );
    if (Platform.OS === 'web') {
      sessionStorage.removeItem('tokenResponse');
    }
  } catch (err) {
    console.error(err);
  }
};
