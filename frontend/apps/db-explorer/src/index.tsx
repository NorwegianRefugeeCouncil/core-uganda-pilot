import React from 'react';
import ReactDOM from 'react-dom';

import './index.css';
import * as log from 'loglevel';

import App from './App';

log.setDefaultLevel(log.levels.TRACE);

ReactDOM.render(
  <React.StrictMode>
    <App/>
  </React.StrictMode>,
  document.getElementById('root'),
);

if (module.hot) {
  module.hot.accept('./App', () => {
    const NextApp = require('./App').default;
    ReactDOM.render(
      <React.StrictMode>
        <NextApp/>
      </React.StrictMode>,
      document.getElementById('root'),
    );
  });
}
