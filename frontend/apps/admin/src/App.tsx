import React from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {BrowserRouter, Switch} from 'react-router-dom';
import {AuthWrapper} from 'core-auth';

import Client from './client/client';
import AuthenticatedApp from './components/AuthenticatedApp';

function App() {
    const client = new Client(
        process?.env?.REACT_APP_AUTHNZ_API_SERVER_URI
    );

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
                    clientId={clientId}
                    issuer={issuer}
                    scopes={scopes}
                    redirectUri={redirectUri}
                    axiosInstance={client.axiosInstance}
                >
                    <AuthenticatedApp/>
                </AuthWrapper>
            </Switch>
        </BrowserRouter>
    );
}

export default App;
