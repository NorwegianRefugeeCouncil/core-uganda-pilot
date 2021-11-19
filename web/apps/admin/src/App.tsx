import React from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {BrowserRouter, Switch} from "react-router-dom";
import AuthenticatedApp from "./components/AuthenticatedApp";
import {AuthWrapper} from "@core/auth-client";
import Client from "../src/client/client";

function App() {
    const client = new Client();
    return (
        <BrowserRouter>
            <Switch>
                <AuthWrapper
                    clientId={process?.env?.REACT_APP_CLIENT_ID || ''}
                    scopes={['openid', 'profile', 'offline_access']}
                    issuer={process?.env?.REACT_APP_ISSUER || ''}
                    redirectUriSuffix={'app'}
                    axiosInstance={client.axiosInstance}
                >
                    <AuthenticatedApp/>
                </AuthWrapper>
            </Switch>
        </BrowserRouter>
    );
}

export default App;
