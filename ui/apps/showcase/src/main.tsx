import { StrictMode } from 'react';
import * as ReactDOM from 'react-dom';

// import { BrowserRouter } from 'react-router-dom';

import App from './app/app';

ReactDOM.render(
  <StrictMode>
    <App />
  </StrictMode>,
  document.getElementById('root')
);
