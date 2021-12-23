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
    if (
      discovery == null ||
      request?.codeVerifier == null ||
      response == null ||
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

    exchangeCodeAsync(exchangeConfig, discovery)
      .then((resp) => {
        setTokenResponse(resp);
      })
      .catch((err) => {
        console.log('Code Exchange Error', err);
        setTokenResponse(undefined);
      });
  }, [request?.codeVerifier, response, discovery]);

  // Trigger onTokenChange callback when TokenResponse is received
  useEffect(() => {
    console.log('tokenResponse has changed...');

    if (tokenResponse == null) return;

    let token;
    switch (injectToken) {
      case 'access_token':
        token = tokenResponse?.accessToken ?? '';
        break;
      case 'id_token':
        token = tokenResponse?.idToken ?? '';
        break;
      default:
        token = '';
    }
    onTokenChange(token);
  }, [JSON.stringify(tokenResponse)]);

  // Update logged in status accordingly
  useEffect(() => {
    if (tokenResponse != null) {
      if (!isLoggedIn) {
        setIsLoggedIn(true);
      }
    } else if (isLoggedIn) {
      setIsLoggedIn(false);
    }
  }, [tokenResponse, isLoggedIn]);

  // Schedule a token refresh before it expires (or use an existing schedule)
  const refreshTokenTimeout = useRef<number | null>(null);
  useEffect(() => {
    if (
      tokenResponse != null &&
      tokenResponse.expiresIn != null &&
      tokenResponse.expiresIn > 0
    ) {
      if (refreshTokenTimeout.current != null)
        window.clearTimeout(refreshTokenTimeout.current);

      const delay = tokenResponse.getExpiryMs() - Date.now();
      console.log({ delay });

      if (delay <= 0) {
        refreshToken();
      } else {
        refreshTokenTimeout.current = window.setTimeout(
          () => refreshToken(),
          delay,
        );
      }
    }
  }, [JSON.stringify(tokenResponse), discovery]);

  const handleLogin = useCallback(() => {
    promptAsync().catch((err) => {
      handleLoginErr(err);
    });
  }, [discovery, request, promptAsync]);

  async function refreshToken() {
    if (tokenResponse != null && discovery != null) {
      const refreshConfig = {
        clientId,
        scopes,
        extraParams: {},
      };
      let response;
      try {
        response = await tokenResponse.refreshAsync(refreshConfig, discovery);
      } catch (err) {
        console.error('Problem encountered while refreshing token', err);
        setTokenResponse(undefined);
      } finally {
        console.log('Refreshed Token: ', response?.accessToken);
        setTokenResponse(response);
      }
    }
  }

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
