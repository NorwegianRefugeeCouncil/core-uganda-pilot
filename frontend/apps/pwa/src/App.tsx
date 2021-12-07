import React, {Fragment} from 'react';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {BrowserRouter, Route, Switch} from "react-router-dom";
import {AuthWrapper} from "core-auth";
import AuthenticatedApp from "./components/AuthenticatedApp";
import client from "./app/client";

const App: React.FC = () => {

    const scope = process?.env?.REACT_APP_OAUTH_SCOPE
    const clientId = process?.env?.REACT_APP_OAUTH_CLIENT_ID;
    const issuer = process?.env?.REACT_APP_OIDC_ISSUER
    const redirectUri = process?.env?.REACT_APP_OAUTH_REDIRECT_URI
    if (!clientId || !issuer || !scope || !redirectUri) {
        return <Fragment/>
    }
    const scopes = scope.split(" ")
    const baseUrl = new URL(redirectUri)
    return (
        <BrowserRouter basename={`${baseUrl.pathname}`}>
            <Switch>
                <Route path={""} render={() =>
                    <AuthWrapper
                        clientId={clientId}
                        axiosInstance={client.axiosInstance}
                        scopes={scopes}
                        issuer={issuer}
                        redirectUri={redirectUri}
                    >
                        <AuthenticatedApp/>
                    </AuthWrapper>
                }/>
            </Switch>
        </BrowserRouter>
    );
}

export default App;

