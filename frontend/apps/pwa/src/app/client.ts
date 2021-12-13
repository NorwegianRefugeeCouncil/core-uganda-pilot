import { Client } from 'core-api-client';
import axios from 'axios';

export default new Client(`${process.env.REACT_APP_SERVER_URL}`, axios.create());
