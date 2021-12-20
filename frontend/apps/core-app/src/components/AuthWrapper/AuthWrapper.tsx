import * as React from 'react';
import { Platform } from 'react-native';
import {
  CodeChallengeMethod,
  exchangeCodeAsync,
  makeRedirectUri,
  ResponseType,
  TokenResponse,
  useAuthRequest,
  useAutoDiscovery,
} from 'expo-auth-session';
import Constants from 'expo-constants';

import { LoginScreen } from '../../screens/LoginScreen';

type Props = {
  onTokenChange: (accessToken: string) => any;
  children: any;
};

export const AuthWrapper: React.FC<Props> = ({ onTokenChange, children }) => {
  const [loggedIn, setLoggedIn] = React.useState(false);
  const [tokenResponse, setTokenResponse] = React.useState<TokenResponse>();

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

  React.useEffect(() => {
    (async () => {
      if (!discovery) return;
      if (!request?.codeVerifier) return;
      if (!response || response.type !== 'success') return;

      const exchangeConfig = {
        code: response.params.code,
        clientId,
        redirectUri,
        extraParams: {
          code_verifier: request?.codeVerifier,
        },
      };

      try {
        const tr = await exchangeCodeAsync(exchangeConfig, discovery);
        setTokenResponse(tr);
      } catch {
        setTokenResponse(undefined);
      }
    })();
  }, [request?.codeVerifier, JSON.stringify(response), JSON.stringify(discovery)]);

  React.useEffect(() => {
    (async () => {
      if (!discovery) return;

      if (tokenResponse?.shouldRefresh()) {
        const refreshConfig = {
          clientId,
          scopes: Constants.manifest?.extra?.scopes,
          extraParams: {},
        };

        try {
          const resp = await tokenResponse?.refreshAsync(refreshConfig, discovery);
          setTokenResponse(resp);
        } catch {
          setTokenResponse(undefined);
        }
      }
    })();
  }, [tokenResponse?.shouldRefresh(), discovery]);

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
