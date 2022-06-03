import {
  Identity,
  IdentityClientDefinition,
  IdentityGetRequest,
  IdentityGetResponse,
} from '../types';

import { BaseRESTClient } from './BaseRESTClient';

export class IdentityClient implements IdentityClientDefinition {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  get = async (request: IdentityGetRequest): Promise<IdentityGetResponse> => {
    return this.restClient.get<Identity>(`/identities/${request.id}`);
  };
}
