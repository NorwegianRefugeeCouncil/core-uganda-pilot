import { RESTClient } from '../rest';
import { ClientSet } from './interfaces';
import fetchMock from 'jest-fetch-mock';
import { APIResource } from '@core2/api-client';

fetchMock.enableMocks();

describe('client', () => {
  it('should work', async () => {

    const client = new RESTClient({
      baseUri: 'http://mocked-by-jest-fetch-mock'
    });

    const cs = new ClientSet(client);

    const res: APIResource = {
      name: 'name',
      kind: 'kind',
      singularName: 'singularName',
      group: 'group',
      namespaced: false,
      verbs: [],
      version: 'v1'
    };
    fetchMock.mockOnce(JSON.stringify(res));

    try {
      const response = await cs.discovery().v1().apiGroups().get('core.nrc.no');
      console.log(response);
      return;
    } catch (e) {
      console.error(e);
      throw e;
    }

  });
});
