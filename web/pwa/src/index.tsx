import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {store} from './app/store';
import {Provider} from 'react-redux';
import * as serviceWorker from './serviceWorker';
import {AuthProvider} from 'oidc-react';
import * as log from "loglevel"
import {SQLContextProvider} from "./app/db";

log.setDefaultLevel(log.levels.TRACE)

const oidcConfig = {
    authority: process.env.REACT_APP_ISSUER,
    clientId: process.env.REACT_APP_CLIENT_ID,
    redirectUri: process.env.REACT_APP_REDIRECT_URI,
    silentRedirectUri: process.env.REACT_APP_SILENT_REDIRECT_URI,
};

ReactDOM.render(
    <React.StrictMode>
        <AuthProvider
            scope={"openid profile email"}
            autoSignIn={true}
            automaticSilentRenew={true}

            {...oidcConfig} >
            <Provider store={store}>
                <SQLContextProvider>
                    <App/>
                </SQLContextProvider>
            </Provider>
        </AuthProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
