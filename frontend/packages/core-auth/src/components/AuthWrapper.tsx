import React, { useCallback, useEffect, useMemo, useState } from 'react';

import Browser from '../types/browser';
import exchangeCodeAsync from '../utils/exchangeCodeAsync';
import useDiscovery from '../hooks/useDiscovery';
import useAuthRequest from '../hooks/useAuthRequest';
import {
  AuthWrapperProps,
  CodeChallengeMethod,
  ResponseType,
  TokenResponseConfig,
} from '../types/types';
import { TokenResponse } from '../types/response';
import { getSessionStorage } from '../utils/getSessionStorage';

const createTokenResponse = (trc: TokenResponseConfig) => {
  if (!trc) return undefined;
  let t: TokenResponse;
  try {
    t = new TokenResponse(trc);
  } catch (e) {
    return undefined;
  }
  return t;
};

const AuthWrapper: React.FC<AuthWrapperProps> = ({
  children,
  scopes = [],
  clientId,
  issuer,
  redirectUri,
  customLoginComponent,
  handleLoginErr = console.log,
  onTokenChange,
  injectToken = 'access_token',
}) => {
  const browser = useMemo(() => new Browser(), []);

  browser.maybeCompleteAuthSession();

  const discovery = useDiscovery(issuer);

  const [tokenResponse, setTokenResponse] = useState<TokenResponse | undefined>(
    createTokenResponse(getSessionStorage(injectToken)),
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

  // Run onTokenChange callback
  // store token in session storage or remove it when undefined
  const applyToken = (tr: TokenResponse | undefined) => {
    if (!tr) {
      sessionStorage.removeItem(injectToken);
    } else {
      const token = (() => {
        switch (injectToken) {
          case 'access_token':
            return tr?.accessToken ?? '';
          case 'id_token':
            return tr?.idToken ?? '';
          default:
            return '';
        }
      })();
      onTokenChange(token);
      sessionStorage.setItem(injectToken, JSON.stringify(tr));
    }
  };
  applyToken(tokenResponse);

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
        applyToken(tr);
        setTokenResponse(tr);
      } catch {
        applyToken(undefined);
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
      };

      try {
        const resp = await tokenResponse.refreshAsync(refreshConfig, discovery);
        applyToken(resp);
        setTokenResponse(resp);
      } catch (err) {
        applyToken(undefined);
        setTokenResponse(undefined);
        setIsLoggedIn(false);
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
  }, [tokenResponse?.refreshToken, tokenResponse?.expiresIn, discovery]);

  // Update logged in status accordingly
  useEffect(() => {
    if (tokenResponse) {
      if (!isLoggedIn) {
        setIsLoggedIn(true);
      }
    } else if (isLoggedIn) {
      setIsLoggedIn(false);
    }
  }, [tokenResponse, isLoggedIn]);

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
            Login {tokenResponse?.accessToken}
          </button>
        )}
      </>
    );
  }

  return <>{children}</>;
};

export default AuthWrapper;
