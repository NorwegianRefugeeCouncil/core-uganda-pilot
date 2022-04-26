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

const getStoredTokenResponse = (): TokenResponse | undefined => {
  try {
    if (Platform.OS !== 'web') return undefined;
    const stored = sessionStorage.getItem('tokenResponse');
    if (stored) {
      return new TokenResponse(JSON.parse(stored));
    }
    return undefined;
  } catch (e) {
    return undefined;
  }
};

const storeTokenResponse = (tokenResponse: TokenResponse) => {
  try {
    if (Platform.OS !== 'web') return;
    if (tokenResponse) {
      sessionStorage.setItem('tokenResponse', JSON.stringify(tokenResponse));
    } else {
      sessionStorage.removeItem('tokenResponse');
    }
  } catch (e) {
    // ignore
  }
};

export const useTokenResponse = (): [
  TokenResponse | undefined,
  () => Promise<void>,
] => {
  const [tokenResponse, setTokenResponse] = React.useState<
    TokenResponse | undefined
  >(getStoredTokenResponse());

  const shouldUseProxy = React.useMemo(
    () => Platform.select({ web: false, default: false }),
    [],
  );

  const redirectUri = React.useMemo(
    () => makeRedirectUri({ scheme: Constants.manifest?.scheme }),
    [],
  );

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
  }, [
    request?.codeVerifier,
    JSON.stringify(response),
    JSON.stringify(discovery),
  ]);

  const refreshTokenInterval = React.useRef<number | null>(null);
  React.useEffect(() => {
    const refreshToken = async () => {
      if (!discovery) return;
      if (!tokenResponse) return;
      if (!tokenResponse.shouldRefresh()) return;

      const refreshConfig = {
        clientId,
        scopes: Constants.manifest?.extra?.scopes,
        extraParams: {},
        refreshToken: tokenResponse.refreshToken,
      };

      try {
        const resp = await tokenResponse.refreshAsync(refreshConfig, discovery);
        setTokenResponse(resp);
      } catch (err) {
        setTokenResponse(undefined);
      }
    };

    if (refreshTokenInterval.current)
      window.clearInterval(refreshTokenInterval.current);

    if (
      tokenResponse &&
      tokenResponse.expiresIn &&
      tokenResponse.expiresIn > 0
    ) {
      refreshTokenInterval.current = window.setInterval(refreshToken, 1000);
    }

    return () => {
      if (refreshTokenInterval.current)
        clearInterval(refreshTokenInterval.current);
    };
  }, [
    tokenResponse?.refreshToken,
    tokenResponse?.expiresIn,
    JSON.stringify(discovery),
  ]);

  React.useEffect(() => {
    if (tokenResponse) storeTokenResponse(tokenResponse);
  }, [JSON.stringify(tokenResponse)]);

  const login = React.useCallback(async () => {
    promptAsync({ useProxy: shouldUseProxy });
  }, [shouldUseProxy, promptAsync]);

  return [tokenResponse, login];
};
