import * as React from 'react';

import { LoginScreen } from '../../screens/LoginScreen';

import { useTokenResponse } from './useTokenResponse';

type Props = {
  onTokenChange: (accessToken: string) => any;
  children: any;
};

export const AuthWrapper: React.FC<Props> = ({ onTokenChange, children }) => {
  const [loggedIn, setLoggedIn] = React.useState(false);

  const [tokenResponse, login] = useTokenResponse();

  React.useEffect(() => {
    if (tokenResponse) {
      if (!loggedIn) {
        setLoggedIn(true);
      }
    } else if (loggedIn) {
      setLoggedIn(false);
    }
  }, [tokenResponse, loggedIn]);

  React.useEffect(() => {
    onTokenChange(tokenResponse?.accessToken ?? '');
  }, [tokenResponse?.accessToken]);

  if (!loggedIn) {
    return <LoginScreen onLogin={login} />;
  }

  return <>{children}</>;
};
