import { ClientSet, RESTClient } from '@core/api-client';

const client = new RESTClient({
  baseUri: 'http://localhost:8000/apis'
});

const Api = new ClientSet(client);

export default Api
