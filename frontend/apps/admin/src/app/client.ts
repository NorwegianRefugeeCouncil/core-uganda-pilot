import axios from 'axios';

import Client from '../client/client';

export const client = new Client(`${process?.env?.REACT_APP_AUTHNZ_API_SERVER_URI}`, axios.create());
