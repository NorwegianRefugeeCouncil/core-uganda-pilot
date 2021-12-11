import axios from 'axios';
import React, { useCallback, useEffect, useState } from 'react';

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
  axiosInstance = axios.create(),
  issuer,
  redirectUri,
  customLoginComponent,
  handleLoginErr = console.log,
  injectToken = 'access_token',
}) => {
  const browser = new Browser();
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
        console.log('Code Exchange Error', err);
        setTokenResponse(undefined);
      });
  }, [request?.codeVerifier, response, discovery]);

  useEffect(() => {
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
        .catch((err) => {
          setTokenResponse(undefined);
        });
    }
  }, [tokenResponse?.shouldRefresh(), discovery]);

  useEffect(() => {
    if (tokenResponse) {
      if (!isLoggedIn) {
        setIsLoggedIn(true);
      }
    } else if (isLoggedIn) {
      setIsLoggedIn(false);
    }
  }, [tokenResponse, isLoggedIn]);

  useEffect(() => {
    const interceptor = axiosInstance.interceptors.request.use((value) => {
      const result = value;
      if (!tokenResponse) {
        return result;
      }
      if (injectToken === 'access_token' && !tokenResponse.accessToken) {
        return result;
      }
      if (injectToken === 'id_token' && !tokenResponse.idToken) {
        return result;
      }
      if (!result.headers) {
        result.headers = {};
      }
      result.headers.Authorization = `Bearer ${
        injectToken === 'access_token' ? tokenResponse?.accessToken : tokenResponse?.idToken
      }`;
      return value;
    });
    return () => {
      axiosInstance.interceptors.request.eject(interceptor);
    };
  }, [tokenResponse?.accessToken]);

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
