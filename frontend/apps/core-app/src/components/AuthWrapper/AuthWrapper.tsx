import * as React from 'react';
import { Platform } from 'react-native';
import { CodeChallengeMethod, makeRedirectUri, ResponseType, useAuthRequest, useAutoDiscovery } from 'expo-auth-session';
import Constants from 'expo-constants';

import { LoginScreen } from '../../screens/LoginScreen';

import { useTokenResponse } from './useTokenResponse';

type Props = {
  onTokenChange: (accessToken: string) => any;
  children: any;
};

export const AuthWrapper: React.FC<Props> = ({ onTokenChange, children }) => {
  const [loggedIn, setLoggedIn] = React.useState(false);

  const shouldUseProxy = React.useMemo(() => Platform.select({ web: false, default: false }), []);
  const redirectUri = React.useMemo(() => makeRedirectUri({ scheme: Constants.manifest?.scheme }), []);

  const discovery = useAutoDiscovery(Constants.manifest?.extra?.issuer);
  const clientId = Constants.manifest?.extra?.client_id;

  const [request, response, promptAsync] = useAuthRequest(
    {
      clientId,
      usePKCE: true,
      responseType: ResponseType.Code,
      codeChallengeMethod: CodeChallengeMethod.S256,
      scopes: Constants.manifest?.extra?.scopes,
      redirectUri,
    },
    discovery,
  );

  const tokenResponse = useTokenResponse(discovery, request, response, clientId, redirectUri);

  React.useEffect(() => {
    if (tokenResponse) {
      if (!loggedIn) {
        setLoggedIn(true);
      }
    } else if (loggedIn) {
      setLoggedIn(false);
    }
  }, [tokenResponse, loggedIn]);

  React.useEffect(() => {
    onTokenChange(tokenResponse?.accessToken ?? '');
  }, [tokenResponse?.accessToken]);

  const handleLogin = React.useCallback(async () => {
    promptAsync({ useProxy: shouldUseProxy });
  }, [shouldUseProxy, promptAsync]);

  if (!loggedIn) {
    return <LoginScreen onLogin={handleLogin} />;
  }

  return <>{children}</>;
};
