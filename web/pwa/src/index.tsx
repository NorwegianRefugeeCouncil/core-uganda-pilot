import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {store} from './app/store';
import {Provider} from 'react-redux';
import * as serviceWorker from './serviceWorker';
import * as log from "loglevel"
import {SQLContextProvider} from "./app/db";

log.setDefaultLevel(log.levels.TRACE)

ReactDOM.render(
    <React.StrictMode>
        <Provider store={store}>
            <SQLContextProvider>
                <App/>
            </SQLContextProvider>
        </Provider>
    </React.StrictMode>,
    document.getElementById('root')
);

// If you want your hooks to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
