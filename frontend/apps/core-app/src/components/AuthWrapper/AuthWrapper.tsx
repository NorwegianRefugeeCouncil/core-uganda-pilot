import * as React from 'react';

import { LoginScreen } from '../../screens/LoginScreen';

import { useTokenResponse } from './useTokenResponse';

type Props = {
  onTokenChange: (accessToken: string) => any;
  children: any;
};

export const AuthWrapper: React.FC<Props> = ({ onTokenChange, children }) => {
  const [loggedIn, setLoggedIn] = React.useState(false);
  const [discoveryLoading, setDiscoveryLoading] = React.useState(true);

  const [tokenResponse, login] = useTokenResponse(setDiscoveryLoading);

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
    return <LoginScreen onLogin={login} isLoading={discoveryLoading} />;
  }

  return <>{children}</>;
};
