import React from 'react';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {BrowserRouter, Route, Switch} from "react-router-dom";
import {AuthWrapper} from "@core/auth-client";
import defaultAuth from "./constants/defaultAuth";
import AuthenticatedApp from "./components/AuthenticatedApp";
import client from "./app/client";

const App: React.FC = () => {
    return (
        <BrowserRouter basename={`/${defaultAuth.redirectUriSuffix}`}>
            <Switch>
                <Route path={""} render={() =>
                    <AuthWrapper
                        clientId={process?.env?.REACT_APP_CLIENT_ID || defaultAuth.clientId}
                        axiosInstance={client.axiosInstance}
                        scopes={defaultAuth.scopes}
                        issuer={process?.env?.REACT_APP_ISSUER || defaultAuth.issuer}
                        redirectUriSuffix={defaultAuth.redirectUriSuffix}
                    >
                        <AuthenticatedApp/>
                    </AuthWrapper>
                }/>
            </Switch>
        </BrowserRouter>
    );
}

export default App;

