import {
  FormCreateRequest,
  FormCreateResponse,
  FormGetRequest,
  FormGetResponse,
  FormListResponse,
} from '../types/client/Form';

import { BaseRESTClient } from './BaseRESTClient';

export class FormClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = (request: FormCreateRequest): Promise<FormCreateResponse> => {
    return this.restClient.do(request, '/forms', 'post', request.object, 200);
  };

  list = (request: {} | undefined): Promise<FormListResponse> => {
    return this.restClient.do(request, '/forms', 'get', undefined, 200);
  };

  get = (request: FormGetRequest): Promise<FormGetResponse> => {
    return this.restClient.do(
      request,
      `/forms/${request.id}`,
      'get',
      undefined,
      200,
    );
  };
}
