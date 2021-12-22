import React, { useCallback, useEffect, useMemo, useState } from 'react';

import Browser from '../types/browser';
import exchangeCodeAsync from '../utils/exchangeCodeAsync';
import useDiscovery from '../hooks/useDiscovery';
import useAuthRequest from '../hooks/useAuthRequest';
import { AuthWrapperProps, CodeChallengeMethod, ResponseType } from '../types/types';
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

  useEffect(() => {
    if (!discovery) {
      return;
    }
    if (!request?.codeVerifier) {
      return;
    }
    if (!response || response?.type !== 'success') {
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

    exchangeCodeAsync(exchangeConfig, discovery)
      .then((resp) => {
        setTokenResponse(resp);
      })
      .catch((err) => {
        setTokenResponse(undefined);
      });
  }, [request?.codeVerifier, response, discovery]);

  const refreshToken = () => {
    if (!discovery) {
      return;
    }
    if (tokenResponse?.shouldRefresh()) {
      const refreshConfig = {
        clientId,
        scopes,
        extraParams: {},
      };
      tokenResponse
        ?.refreshAsync(refreshConfig, discovery)
        .then((resp) => {
          setTokenResponse(resp);
        })
        .catch(() => {
          setTokenResponse(undefined);
        });
    }
  };
  let refreshTokenIntervalId: NodeJS.Timeout | null = null;

  useEffect(() => {
    if (tokenResponse) {
      if (!isLoggedIn) {
        setIsLoggedIn(true);
      }
      if (refreshTokenIntervalId != null) clearInterval(refreshTokenIntervalId);
      refreshTokenIntervalId = setInterval(()=>refreshToken(), 5000);
    } else if (isLoggedIn) {
      setIsLoggedIn(false);
    }
  }, [tokenResponse, isLoggedIn]);

  useEffect(() => {
    onTokenChange(token);
  }, [token]);

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
