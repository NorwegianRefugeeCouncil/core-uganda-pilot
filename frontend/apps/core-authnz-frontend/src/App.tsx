import React from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import { AuthWrapper } from 'core-auth';

import { client } from './app/client';
import { Router } from './components/Router';

const App: React.FC = () => {
  const scope = process?.env?.REACT_APP_OAUTH_SCOPE;
  const clientId = process?.env?.REACT_APP_OAUTH_CLIENT_ID;
  const issuer = process?.env?.REACT_APP_OIDC_ISSUER;
  const redirectUri = process?.env?.REACT_APP_OAUTH_REDIRECT_URI;
  if (!clientId || !issuer || !scope || !redirectUri) {
    return <></>;
  }
  const scopes = scope.split(' ');
  const baseUrl = new URL(redirectUri);

  return (
    <AuthWrapper
      injectToken="id_token"
      onTokenChange={client.setToken}
      clientId={clientId}
      issuer={issuer}
      scopes={scopes}
      redirectUri={redirectUri}
    >
      <Router baseUrl={baseUrl.pathname} />
    </AuthWrapper>
  );
};

export default App;
