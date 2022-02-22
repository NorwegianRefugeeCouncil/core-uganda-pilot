import {
  RecordCreateRequest,
  RecordCreateResponse,
  RecordGetRequest,
  RecordGetResponse,
  RecordListRequest,
  RecordListResponse,
} from 'core-api-client';
import { RecordClient } from 'core-api-client/src/lib/client/Record';

export class RecipientsClient {
  restClient: RecordClient;

  constructor(restClient: RecordClient) {
    this.restClient = restClient;
  }

  create = (request: RecordCreateRequest): Promise<RecordCreateResponse> => {
    return this.restClient.create(request);
  };

  list = (request: RecordListRequest): Promise<RecordListResponse> => {
    return this.restClient.list(request);
  };

  get = (request: RecordGetRequest): Promise<RecordGetResponse> => {
    return this.restClient.get(request);
  };
}
