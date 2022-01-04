import React, {
  useCallback,
  useEffect,
  useMemo,
  useState,
  useRef,
} from 'react';

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

  const [tokenResponse, setTokenResponse] = useState<TokenResponse>();

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
  useEffect(() => {
    (async () => {
      if (
        !discovery ||
        request?.codeVerifier == null ||
        !response ||
        response?.type !== 'success'
      )
        return;

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

  // Trigger onTokenChange callback when TokenResponse is received
  useEffect(() => {
    if (!tokenResponse) return;

    let token;
    switch (injectToken) {
      case 'access_token':
        token = tokenResponse?.accessToken || '';
        break;
      case 'id_token':
        token = tokenResponse?.idToken || '';
        break;
      default:
        token = '';
    }
    onTokenChange(token);
  }, [tokenResponse?.accessToken, tokenResponse?.idToken]);

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

  async function refreshToken() {
    if (tokenResponse && discovery) {
      const refreshConfig = {
        clientId,
        scopes,
        extraParams: {},
      };
      let newTokenResponse;
      try {
        newTokenResponse = await tokenResponse.refreshAsync(
          refreshConfig,
          discovery,
        );
      } catch (err) {
        setTokenResponse(undefined);
        throw err;
      }

      setTokenResponse(newTokenResponse);
    }
  }

  // Schedule a token refresh before it expires (or use an existing schedule)
  const refreshTokenInterval = useRef<number | null>(null);
  useEffect(() => {
    if (
      tokenResponse &&
      tokenResponse.expiresIn != null &&
      tokenResponse.expiresIn > 0
    ) {
      if (refreshTokenInterval.current != null)
        window.clearInterval(refreshTokenInterval.current);

      const delay = tokenResponse.getExpiryMs() - Date.now();
      if (delay <= 0) {
        refreshToken();
      } else {
        refreshTokenInterval.current = window.setInterval(
          () => refreshToken(),
          delay / 3,
        );
      }
    }
  }, [tokenResponse?.accessToken, discovery]);

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
