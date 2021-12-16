import React from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import { BrowserRouter, Switch } from 'react-router-dom';
import { AuthWrapper } from 'core-auth';

import AuthenticatedApp from './components/AuthenticatedApp';
import { client } from './app/client';

function App() {
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
    <BrowserRouter basename={baseUrl.pathname}>
      <Switch>
        <AuthWrapper
          injectToken="id_token"
          axiosInstance={client.axiosInstance}
          clientId={clientId}
          issuer={issuer}
          scopes={scopes}
          redirectUri={redirectUri}
        >
          <AuthenticatedApp />
        </AuthWrapper>
      </Switch>
    </BrowserRouter>
  );
}

export default App;
