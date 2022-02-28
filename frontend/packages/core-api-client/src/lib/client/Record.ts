import {
  RecordCreateRequest,
  RecordCreateResponse,
  RecordGetRequest,
  RecordGetResponse,
  RecordListRequest,
  RecordListResponse,
} from '../types';

import { BaseRESTClient } from './BaseRESTClient';

export class RecordClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = (request: RecordCreateRequest): Promise<RecordCreateResponse> => {
    return this.restClient.do(request, '/records', 'post', request.object, 200);
  };

  list = (request: RecordListRequest): Promise<RecordListResponse> => {
    const { databaseId, formId } = request;
    const url = `/records?databaseId=${databaseId}&formId=${formId}`;
    return this.restClient.do(request, url, 'get', undefined, 200);
  };

  get = (request: RecordGetRequest): Promise<RecordGetResponse> => {
    const { databaseId, formId, recordId } = request;
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    return this.restClient.do(request, url, 'get', undefined, 200);
  };
}
