import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { getLogger, LogLevelDesc } from 'loglevel';

import Browser from '../types/browser';
import exchangeCodeAsync from '../utils/exchangeCodeAsync';
import useDiscovery from '../hooks/useDiscovery';
import useAuthRequest from '../hooks/useAuthRequest';
import {
  AuthWrapperProps,
  CodeChallengeMethod,
  ResponseType,
} from '../types/types';
import { TokenResponse } from '../types/response';
import { getJSONFromSessionStorage } from '../utils/getJSONFromSessionStorage';

const log = getLogger('AuthWrapper');
const loglevel = process?.env?.LOG_LEVEL as LogLevelDesc;
log.setLevel(loglevel || log.levels.INFO);

const AuthWrapper: React.FC<AuthWrapperProps> = ({
  children,
  scopes = [],
  clientId,
  issuer,
  redirectUri,
  customLoginComponent,
  handleLoginErr = log.debug,
  onTokenChange,
  injectToken = 'access_token',
}) => {
  const browser = useMemo(() => new Browser(), []);

  browser.maybeCompleteAuthSession();

  const discovery = useDiscovery(issuer);

  const [tokenResponse, setTokenResponse] = useState<TokenResponse | undefined>(
    TokenResponse.createTokenResponse(getJSONFromSessionStorage(injectToken)),
  );

  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const [request, response, promptAsync] = useAuthRequest(
    {
      clientId,
      usePKCE: true,
      responseType: ResponseType.Code,
      codeChallengeMethod: CodeChallengeMethod.S256,
      scopes,
      redirectUri,
    },
    discovery,
    browser,
  );

  // Make initial token request
  // trigger login automatically
  useEffect(() => {
    (async () => {
      if (
        !discovery ||
        request?.codeVerifier == null ||
        !response ||
        response?.type !== 'success'
      ) {
        return;
      }

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
  }, [request?.codeVerifier, response, discovery]);

  // Each second, check if token is fresh
  // If not, refresh the token
  const refreshTokenInterval = React.useRef<number | null>(null);
  React.useEffect(() => {
    const refreshToken = async () => {
      if (!discovery) return;
      if (!tokenResponse) return;
      if (!tokenResponse.shouldRefresh()) return;

      const refreshConfig = {
        clientId,
        scopes,
        extraParams: {},
        refreshToken: tokenResponse.refreshToken,
      };

      try {
        const resp = await TokenResponse.refreshAsync(refreshConfig, discovery);
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

  // Run onTokenChange callback
  // store token in session storage or remove it when undefined
  useEffect(() => {
    if (!tokenResponse) {
      onTokenChange('');
      sessionStorage.removeItem(injectToken);
    } else {
      const token = (() => {
        switch (injectToken) {
          case 'access_token':
            return tokenResponse?.accessToken ?? '';
          case 'id_token':
            return tokenResponse?.idToken ?? '';
          default:
            return '';
        }
      })();
      onTokenChange(token);
      sessionStorage.setItem(injectToken, JSON.stringify(tokenResponse));
    }
  }, [tokenResponse?.accessToken]);

  // Update logged in status accordingly
  useEffect(() => {
    if (tokenResponse) {
      if (!isLoggedIn) {
        setIsLoggedIn(true);
      }
    } else if (isLoggedIn) {
      setIsLoggedIn(false);
    }
  }, [JSON.stringify(tokenResponse), isLoggedIn]);

  // trigger login manually, by clicking
  const handleLogin = useCallback(() => {
    promptAsync().catch((err) => {
      handleLoginErr(err);
    });
  }, [discovery, request, promptAsync]);

  if (!isLoggedIn) {
    return (
      <>
        {customLoginComponent ? (
          customLoginComponent({ login: handleLogin })
        ) : (
          <button type="button" onClick={handleLogin}>
            Login
          </button>
        )}
      </>
    );
  }

  return <>{children}</>;
};

export default AuthWrapper;
