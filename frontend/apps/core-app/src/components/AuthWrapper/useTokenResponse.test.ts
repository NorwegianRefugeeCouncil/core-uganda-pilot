import { renderHook, act } from '@testing-library/react-hooks';
import { waitFor } from '@testing-library/react-native';
import * as expoAuthSession from 'expo-auth-session';

import { useTokenResponse } from './useTokenResponse';

jest.mock('expo-auth-session', () => {
  return {
    ...jest.requireActual('expo-auth-session'),
    useAutoDiscovery: jest.fn(),
    useAuthRequest: jest.fn(),
    exchangeCodeAsync: jest.fn(),
  };
});

describe('Success', () => {
  it('should call the auth services with the correct parameters', async () => {
    const fakeTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'fakeToken',
    });
    const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
      tokenEndpoint: 'fakeTokenEndpoint',
    };
    const fakeRequest = new expoAuthSession.AuthRequest({
      clientId: 'fake-client-id',
      redirectUri: 'no.nrc.core/redirect',
    });
    fakeRequest.codeVerifier = 'fake-code-verifier';
    const fakeResponse: expoAuthSession.AuthSessionResult = {
      type: 'success',
      errorCode: null,
      params: {},
      authentication: fakeTokenResponse,
      url: '',
    };
    const promptAsyncMock = jest.fn();

    const useAutoDiscoverySpy = jest
      .spyOn(expoAuthSession, 'useAutoDiscovery')
      .mockReturnValue(fakeDiscovery);
    useAutoDiscoverySpy.mockClear();

    const useAuthRequestSpy = jest
      .spyOn(expoAuthSession, 'useAuthRequest')
      .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
    useAuthRequestSpy.mockClear();

    const exchangeCodeAsyncSpy = jest
      .spyOn(expoAuthSession, 'exchangeCodeAsync')
      .mockResolvedValue(fakeTokenResponse);
    exchangeCodeAsyncSpy.mockClear();

    const { waitForNextUpdate } = renderHook(() => useTokenResponse());

    expect(useAutoDiscoverySpy).toHaveBeenCalledTimes(1);
    // issuer defined in core/frontend/app.json
    expect(useAutoDiscoverySpy).toHaveBeenCalledWith('fake-issuer');

    expect(useAuthRequestSpy).toHaveBeenCalledTimes(1);
    expect(useAuthRequestSpy).toHaveBeenCalledWith(
      {
        clientId: 'fake-client-id',
        usePKCE: true,
        responseType: expoAuthSession.ResponseType.Code,
        codeChallengeMethod: expoAuthSession.CodeChallengeMethod.S256,
        scopes: ['fake-scope'],
        redirectUri: 'nrccore://exp.host/@test/frontend-root',
      },
      fakeDiscovery,
    );

    expect(exchangeCodeAsyncSpy).toHaveBeenCalledTimes(1);
    expect(exchangeCodeAsyncSpy).toHaveBeenCalledWith(
      {
        code: fakeResponse.params.code,
        clientId: 'fake-client-id',
        redirectUri: 'nrccore://exp.host/@test/frontend-root',
        extraParams: {
          code_verifier: fakeRequest.codeVerifier,
        },
      },
      fakeDiscovery,
    );

    await waitForNextUpdate();
  });

  it('should return a token response', async () => {
    const fakeTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'fakeToken',
    });
    const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
      tokenEndpoint: 'fakeTokenEndpoint',
    };
    const fakeRequest = new expoAuthSession.AuthRequest({
      clientId: 'fake-client-id',
      redirectUri: 'no.nrc.core/redirect',
    });
    fakeRequest.codeVerifier = 'fake-code-verifier';
    const fakeResponse: expoAuthSession.AuthSessionResult = {
      type: 'success',
      errorCode: null,
      params: {},
      authentication: fakeTokenResponse,
      url: '',
    };
    const promptAsyncMock = jest.fn();

    const useAutoDiscoverySpy = jest
      .spyOn(expoAuthSession, 'useAutoDiscovery')
      .mockReturnValue(fakeDiscovery);
    useAutoDiscoverySpy.mockClear();

    const useAuthRequestSpy = jest
      .spyOn(expoAuthSession, 'useAuthRequest')
      .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
    useAuthRequestSpy.mockClear();

    const exchangeCodeAsyncSpy = jest
      .spyOn(expoAuthSession, 'exchangeCodeAsync')
      .mockResolvedValue(fakeTokenResponse);
    exchangeCodeAsyncSpy.mockClear();

    const { result } = renderHook(() => useTokenResponse());
    await waitFor(() => expect(result.current[0]).toEqual(fakeTokenResponse));
  });

  it('should prompt the user when calling login', async () => {
    const fakeTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'fakeToken',
    });
    const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
      tokenEndpoint: 'fakeTokenEndpoint',
    };
    const fakeRequest = new expoAuthSession.AuthRequest({
      clientId: 'fake-client-id',
      redirectUri: 'no.nrc.core/redirect',
    });
    fakeRequest.codeVerifier = 'fake-code-verifier';
    const fakeResponse: expoAuthSession.AuthSessionResult = {
      type: 'success',
      errorCode: null,
      params: {},
      authentication: fakeTokenResponse,
      url: '',
    };
    const promptAsyncMock = jest.fn();

    const useAutoDiscoverySpy = jest
      .spyOn(expoAuthSession, 'useAutoDiscovery')
      .mockReturnValue(fakeDiscovery);
    useAutoDiscoverySpy.mockClear();

    const useAuthRequestSpy = jest
      .spyOn(expoAuthSession, 'useAuthRequest')
      .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
    useAuthRequestSpy.mockClear();

    const exchangeCodeAsyncSpy = jest
      .spyOn(expoAuthSession, 'exchangeCodeAsync')
      .mockResolvedValue(fakeTokenResponse);
    exchangeCodeAsyncSpy.mockClear();

    const { result, waitForNextUpdate } = renderHook(() => useTokenResponse());

    result.current[1]();
    expect(promptAsyncMock).toHaveBeenCalledTimes(1);
    expect(promptAsyncMock).toHaveBeenCalledWith({ useProxy: false });

    await waitForNextUpdate();
  });

  it('should refresh the token', async () => {
    const fakeRefreshTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'newFakeToken',
    });
    const refreshAsyncMock = jest.fn(() =>
      Promise.resolve(fakeRefreshTokenResponse),
    );
    const fakeTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'fakeToken',
    });
    fakeTokenResponse.shouldRefresh = jest.fn(() => true);
    fakeTokenResponse.refreshAsync = refreshAsyncMock;

    const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
      tokenEndpoint: 'fakeTokenEndpoint',
    };
    const fakeRequest = new expoAuthSession.AuthRequest({
      clientId: 'fake-client-id',
      redirectUri: 'no.nrc.core/redirect',
    });
    fakeRequest.codeVerifier = 'fake-code-verifier';
    const fakeResponse: expoAuthSession.AuthSessionResult = {
      type: 'success',
      errorCode: null,
      params: {},
      authentication: fakeTokenResponse,
      url: '',
    };
    const promptAsyncMock = jest.fn();

    const useAutoDiscoverySpy = jest
      .spyOn(expoAuthSession, 'useAutoDiscovery')
      .mockReturnValue(fakeDiscovery);
    useAutoDiscoverySpy.mockClear();

    const useAuthRequestSpy = jest
      .spyOn(expoAuthSession, 'useAuthRequest')
      .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
    useAuthRequestSpy.mockClear();

    const exchangeCodeAsyncSpy = jest
      .spyOn(expoAuthSession, 'exchangeCodeAsync')
      .mockResolvedValue(fakeTokenResponse);
    exchangeCodeAsyncSpy.mockClear();

    const { result, waitForNextUpdate } = renderHook(() => useTokenResponse());

    expect(useAutoDiscoverySpy).toHaveBeenCalledTimes(1);
    // issuer defined in core/frontend/app.json
    expect(useAutoDiscoverySpy).toHaveBeenCalledWith('fake-issuer');

    expect(useAuthRequestSpy).toHaveBeenCalledTimes(1);
    expect(useAuthRequestSpy).toHaveBeenCalledWith(
      {
        clientId: 'fake-client-id',
        usePKCE: true,
        responseType: expoAuthSession.ResponseType.Code,
        codeChallengeMethod: expoAuthSession.CodeChallengeMethod.S256,
        scopes: ['fake-scope'],
        redirectUri: 'nrccore://exp.host/@test/frontend-root',
      },
      fakeDiscovery,
    );

    expect(exchangeCodeAsyncSpy).toHaveBeenCalledTimes(1);
    expect(exchangeCodeAsyncSpy).toHaveBeenCalledWith(
      {
        code: fakeResponse.params.code,
        clientId: 'fake-client-id',
        redirectUri: 'nrccore://exp.host/@test/frontend-root',
        extraParams: {
          code_verifier: fakeRequest.codeVerifier,
        },
      },
      fakeDiscovery,
    );

    await waitForNextUpdate();

    expect(fakeTokenResponse.refreshAsync).toHaveBeenCalledTimes(1);
    expect(fakeTokenResponse.refreshAsync).toHaveBeenCalledWith(
      {
        clientId: 'fake-client-id',
        scopes: ['fake-scope'],
        extraParams: {},
      },
      fakeDiscovery,
    );
    expect(result.current[0]).toEqual(fakeRefreshTokenResponse);
  });
});

describe('Failure', () => {
  describe('exchangeCodeAsync', () => {
    it('should not exchange codes if the auth request response fails', () => {
      const fakeTokenResponse = new expoAuthSession.TokenResponse({
        accessToken: 'fakeToken',
      });
      const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
        tokenEndpoint: 'fakeTokenEndpoint',
      };
      const fakeRequest = new expoAuthSession.AuthRequest({
        clientId: 'fake-client-id',
        redirectUri: 'no.nrc.core/redirect',
      });
      fakeRequest.codeVerifier = 'fake-code-verifier';
      const fakeResponse: expoAuthSession.AuthSessionResult = {
        type: 'error',
        errorCode: null,
        params: {},
        authentication: fakeTokenResponse,
        url: '',
      };
      const promptAsyncMock = jest.fn();

      const useAutoDiscoverySpy = jest
        .spyOn(expoAuthSession, 'useAutoDiscovery')
        .mockReturnValue(fakeDiscovery);
      useAutoDiscoverySpy.mockClear();

      const useAuthRequestSpy = jest
        .spyOn(expoAuthSession, 'useAuthRequest')
        .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
      useAuthRequestSpy.mockClear();

      const exchangeCodeAsyncSpy = jest
        .spyOn(expoAuthSession, 'exchangeCodeAsync')
        .mockResolvedValue(fakeTokenResponse);
      exchangeCodeAsyncSpy.mockClear();

      const { result } = renderHook(() => useTokenResponse());

      expect(useAutoDiscoverySpy).toHaveBeenCalledTimes(1);
      // issuer defined in core/frontend/app.json
      expect(useAutoDiscoverySpy).toHaveBeenCalledWith('fake-issuer');

      expect(useAuthRequestSpy).toHaveBeenCalledTimes(1);
      expect(useAuthRequestSpy).toHaveBeenCalledWith(
        {
          clientId: 'fake-client-id',
          usePKCE: true,
          responseType: expoAuthSession.ResponseType.Code,
          codeChallengeMethod: expoAuthSession.CodeChallengeMethod.S256,
          scopes: ['fake-scope'],
          redirectUri: 'nrccore://exp.host/@test/frontend-root',
        },
        fakeDiscovery,
      );

      expect(exchangeCodeAsyncSpy).not.toHaveBeenCalled();

      expect(result.current[0]).toBe(undefined);
    });

    it('should not exchange codes if there is no discovery document', () => {
      const fakeTokenResponse = new expoAuthSession.TokenResponse({
        accessToken: 'fakeToken',
      });
      const fakeDiscovery = null;
      const fakeRequest = new expoAuthSession.AuthRequest({
        clientId: 'fake-client-id',
        redirectUri: 'no.nrc.core/redirect',
      });
      fakeRequest.codeVerifier = 'fake-code-verifier';
      const fakeResponse: expoAuthSession.AuthSessionResult = {
        type: 'success',
        errorCode: null,
        params: {},
        authentication: fakeTokenResponse,
        url: '',
      };
      const promptAsyncMock = jest.fn();

      const useAutoDiscoverySpy = jest
        .spyOn(expoAuthSession, 'useAutoDiscovery')
        .mockReturnValue(fakeDiscovery);
      useAutoDiscoverySpy.mockClear();

      const useAuthRequestSpy = jest
        .spyOn(expoAuthSession, 'useAuthRequest')
        .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
      useAuthRequestSpy.mockClear();

      const exchangeCodeAsyncSpy = jest
        .spyOn(expoAuthSession, 'exchangeCodeAsync')
        .mockResolvedValue(fakeTokenResponse);
      exchangeCodeAsyncSpy.mockClear();

      const { result } = renderHook(() => useTokenResponse());

      expect(useAutoDiscoverySpy).toHaveBeenCalledTimes(1);
      // issuer defined in core/frontend/app.json
      expect(useAutoDiscoverySpy).toHaveBeenCalledWith('fake-issuer');

      expect(useAuthRequestSpy).toHaveBeenCalledTimes(1);
      expect(useAuthRequestSpy).toHaveBeenCalledWith(
        {
          clientId: 'fake-client-id',
          usePKCE: true,
          responseType: expoAuthSession.ResponseType.Code,
          codeChallengeMethod: expoAuthSession.CodeChallengeMethod.S256,
          scopes: ['fake-scope'],
          redirectUri: 'nrccore://exp.host/@test/frontend-root',
        },
        fakeDiscovery,
      );

      expect(exchangeCodeAsyncSpy).not.toHaveBeenCalled();

      expect(result.current[0]).toBe(undefined);
    });

    it('should not set the token respone if exchange codes fails', () => {
      const fakeTokenResponse = new expoAuthSession.TokenResponse({
        accessToken: 'fakeToken',
      });
      const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
        tokenEndpoint: 'fakeTokenEndpoint',
      };
      const fakeRequest = new expoAuthSession.AuthRequest({
        clientId: 'fake-client-id',
        redirectUri: 'no.nrc.core/redirect',
      });
      fakeRequest.codeVerifier = 'fake-code-verifier';
      const fakeResponse: expoAuthSession.AuthSessionResult = {
        type: 'success',
        errorCode: null,
        params: {},
        authentication: fakeTokenResponse,
        url: '',
      };
      const promptAsyncMock = jest.fn();

      const useAutoDiscoverySpy = jest
        .spyOn(expoAuthSession, 'useAutoDiscovery')
        .mockReturnValue(fakeDiscovery);
      useAutoDiscoverySpy.mockClear();

      const useAuthRequestSpy = jest
        .spyOn(expoAuthSession, 'useAuthRequest')
        .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
      useAuthRequestSpy.mockClear();

      const exchangeCodeAsyncSpy = jest
        .spyOn(expoAuthSession, 'exchangeCodeAsync')
        .mockImplementation(() => {
          throw new Error('fake error');
        });
      exchangeCodeAsyncSpy.mockClear();

      const { result } = renderHook(() => useTokenResponse());

      expect(useAutoDiscoverySpy).toHaveBeenCalledTimes(1);
      // issuer defined in core/frontend/app.json
      expect(useAutoDiscoverySpy).toHaveBeenCalledWith('fake-issuer');

      expect(useAuthRequestSpy).toHaveBeenCalledTimes(1);
      expect(useAuthRequestSpy).toHaveBeenCalledWith(
        {
          clientId: 'fake-client-id',
          usePKCE: true,
          responseType: expoAuthSession.ResponseType.Code,
          codeChallengeMethod: expoAuthSession.CodeChallengeMethod.S256,
          scopes: ['fake-scope'],
          redirectUri: 'nrccore://exp.host/@test/frontend-root',
        },
        fakeDiscovery,
      );

      expect(exchangeCodeAsyncSpy).toHaveBeenCalledTimes(1);
      expect(exchangeCodeAsyncSpy).toHaveBeenCalledWith(
        {
          code: fakeResponse.params.code,
          clientId: 'fake-client-id',
          redirectUri: 'nrccore://exp.host/@test/frontend-root',
          extraParams: {
            code_verifier: fakeRequest.codeVerifier,
          },
        },
        fakeDiscovery,
      );

      expect(result.current[0]).toBe(undefined);
    });
  });

  it('should fail to refresh the token', async () => {
    const fakeRefreshTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'newFakeToken',
    });
    const refreshAsyncMock = jest.fn(() => {
      throw new Error('fake-error');
    });
    const fakeTokenResponse = new expoAuthSession.TokenResponse({
      accessToken: 'fakeToken',
    });
    fakeTokenResponse.shouldRefresh = jest.fn(() => true);
    fakeTokenResponse.refreshAsync = refreshAsyncMock;

    const fakeDiscovery: expoAuthSession.DiscoveryDocument = {
      tokenEndpoint: 'fakeTokenEndpoint',
    };
    const fakeRequest = new expoAuthSession.AuthRequest({
      clientId: 'fake-client-id',
      redirectUri: 'no.nrc.core/redirect',
    });
    fakeRequest.codeVerifier = 'fake-code-verifier';
    const fakeResponse: expoAuthSession.AuthSessionResult = {
      type: 'success',
      errorCode: null,
      params: {},
      authentication: fakeTokenResponse,
      url: '',
    };
    const promptAsyncMock = jest.fn();

    const useAutoDiscoverySpy = jest
      .spyOn(expoAuthSession, 'useAutoDiscovery')
      .mockReturnValue(fakeDiscovery);
    useAutoDiscoverySpy.mockClear();

    const useAuthRequestSpy = jest
      .spyOn(expoAuthSession, 'useAuthRequest')
      .mockReturnValue([fakeRequest, fakeResponse, promptAsyncMock]);
    useAuthRequestSpy.mockClear();

    const exchangeCodeAsyncSpy = jest
      .spyOn(expoAuthSession, 'exchangeCodeAsync')
      .mockResolvedValue(fakeTokenResponse);
    exchangeCodeAsyncSpy.mockClear();

    const { result, waitForNextUpdate } = renderHook(() => useTokenResponse());

    expect(useAutoDiscoverySpy).toHaveBeenCalledTimes(1);
    // issuer defined in core/frontend/app.json
    expect(useAutoDiscoverySpy).toHaveBeenCalledWith('fake-issuer');

    expect(useAuthRequestSpy).toHaveBeenCalledTimes(1);
    expect(useAuthRequestSpy).toHaveBeenCalledWith(
      {
        clientId: 'fake-client-id',
        usePKCE: true,
        responseType: expoAuthSession.ResponseType.Code,
        codeChallengeMethod: expoAuthSession.CodeChallengeMethod.S256,
        scopes: ['fake-scope'],
        redirectUri: 'nrccore://exp.host/@test/frontend-root',
      },
      fakeDiscovery,
    );

    expect(exchangeCodeAsyncSpy).toHaveBeenCalledTimes(1);
    expect(exchangeCodeAsyncSpy).toHaveBeenCalledWith(
      {
        code: fakeResponse.params.code,
        clientId: 'fake-client-id',
        redirectUri: 'nrccore://exp.host/@test/frontend-root',
        extraParams: {
          code_verifier: fakeRequest.codeVerifier,
        },
      },
      fakeDiscovery,
    );

    await waitForNextUpdate();

    expect(fakeTokenResponse.refreshAsync).toHaveBeenCalledTimes(1);
    expect(fakeTokenResponse.refreshAsync).toHaveBeenCalledWith(
      {
        clientId: 'fake-client-id',
        scopes: ['fake-scope'],
        extraParams: {},
      },
      fakeDiscovery,
    );
    expect(result.current[0]).toEqual(undefined);
  });
});
