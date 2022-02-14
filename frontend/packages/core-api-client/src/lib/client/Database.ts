import {
  DatabaseCreateRequest,
  DatabaseCreateResponse,
  DatabaseListResponse,
} from '../types/client/Database';

import { BaseRESTClient } from './BaseRESTClient';

export class DatabaseClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = (
    request: DatabaseCreateRequest,
  ): Promise<DatabaseCreateResponse> => {
    return this.restClient.do(
      request,
      '/databases',
      'post',
      request.object,
      200,
    );
  };

  list = (request: {} | undefined): Promise<DatabaseListResponse> => {
    return this.restClient.do(request, '/databases', 'get', undefined, 200);
  };
}
