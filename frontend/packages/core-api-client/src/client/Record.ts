import {
  Record,
  RecordCreateRequest,
  RecordCreateResponse,
  RecordGetRequest,
  RecordGetResponse,
  RecordListRequest,
  RecordListResponse,
  RecordList,
} from '../types';

import { BaseRESTClient } from './BaseRESTClient';
import { FormClient } from './Form';

export class RecordClient {
  restClient: BaseRESTClient;

  formClient: FormClient;

  constructor(restClient: BaseRESTClient, formClient: FormClient) {
    this.restClient = restClient;
    this.formClient = formClient;
  }

  create = (request: RecordCreateRequest): Promise<RecordCreateResponse> => {
    return this.restClient.do(request, '/records', 'post', request.object, 200);
  };

  list = (request: RecordListRequest): Promise<RecordListResponse> => {
    const { databaseId, formId } = request;
    const url = `/records?databaseId=${databaseId}&formId=${formId}`;
    return this.restClient.do<RecordListRequest, RecordList>(
      request,
      url,
      'get',
      undefined,
      200,
    );
  };

  get = (request: RecordGetRequest): Promise<RecordGetResponse> => {
    const { databaseId, formId, recordId } = request;
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    return this.restClient.do<RecordGetRequest, Record>(
      request,
      url,
      'get',
      undefined,
      200,
    );
  };
}
