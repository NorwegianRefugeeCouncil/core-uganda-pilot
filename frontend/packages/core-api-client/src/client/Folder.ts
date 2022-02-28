import {
  FolderCreateRequest,
  FolderCreateResponse,
  FolderListResponse,
} from '../types/client/Folder';

import { BaseRESTClient } from './BaseRESTClient';

export class FolderClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = (request: FolderCreateRequest): Promise<FolderCreateResponse> => {
    return this.restClient.do(request, '/folders', 'post', request.object, 200);
  };

  list = (request: {} | undefined): Promise<FolderListResponse> => {
    return this.restClient.do(request, '/folders', 'get', undefined, 200);
  };
}
