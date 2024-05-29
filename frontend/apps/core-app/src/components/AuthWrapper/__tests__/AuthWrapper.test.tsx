import { Text } from 'react-native';
import { fireEvent } from '@testing-library/react-native';
import { TokenResponse } from 'expo-auth-session';
import React, { useState } from 'react';

import { render } from '../../../testUtils/render';
import * as hooks from '../useTokenResponse';
import { AuthWrapper } from '../AuthWrapper';

describe('Unauthenticated', () => {
  it('should render the login button', () => {
    const onTokenChangeMock = jest.fn();
    const loginSpy = jest.fn();

    jest
      .spyOn(hooks, 'useTokenResponse')
      .mockImplementationOnce(() => [undefined, loginSpy]);

    const { getByTestId } = render(
      <AuthWrapper onTokenChange={onTokenChangeMock}>
        <Text>Success</Text>
      </AuthWrapper>,
    );

    const button = getByTestId('login-button');
    expect(button).toBeTruthy();
    fireEvent.press(button);
  });
});

describe('Authenticated', () => {
  it('should call "onTokenChange" when loggedIn', () => {
    const onTokenChangeMock = jest.fn();
    const loginSpy = jest.fn();

    const fakeTokenResponse = new TokenResponse({ accessToken: 'fakeToken' });
    jest
      .spyOn(hooks, 'useTokenResponse')
      .mockImplementation(() => [fakeTokenResponse, loginSpy]);

    render(
      <AuthWrapper onTokenChange={onTokenChangeMock}>
        <Text>Success</Text>
      </AuthWrapper>,
    );

    expect(onTokenChangeMock).toHaveBeenCalledWith('fakeToken');
  });

  it('should render the children', () => {
    const onTokenChangeMock = jest.fn();
    const loginSpy = jest.fn();

    const fakeTokenResponse = new TokenResponse({ accessToken: 'fakeToken' });
    jest
      .spyOn(hooks, 'useTokenResponse')
      .mockImplementation(() => [fakeTokenResponse, loginSpy]);

    const { getByText } = render(
      <AuthWrapper onTokenChange={onTokenChangeMock}>
        <Text>Success</Text>
      </AuthWrapper>,
    );

    expect(getByText('Success')).toBeTruthy();
  });
});
