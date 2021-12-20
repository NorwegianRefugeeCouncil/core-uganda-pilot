import * as React from 'react';
import { AuthRequest, AuthSessionResult, DiscoveryDocument, exchangeCodeAsync, TokenResponse } from 'expo-auth-session';
import Constants from 'expo-constants';

export const useTokenResponse = (
  discovery: DiscoveryDocument | null,
  request: AuthRequest | null,
  response: AuthSessionResult | null,
  clientId: string,
  redirectUri: string,
): TokenResponse | undefined => {
  const [tokenResponse, setTokenResponse] = React.useState<TokenResponse>();

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

  return tokenResponse;
};
