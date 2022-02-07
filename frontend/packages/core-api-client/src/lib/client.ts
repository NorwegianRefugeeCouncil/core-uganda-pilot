import { BaseRESTClient } from './BaseRESTClient';
import {
  ClientDefinition,
  DatabaseCreateRequest,
  DatabaseCreateResponse,
  DatabaseListResponse,
  FolderCreateRequest,
  FolderCreateResponse,
  FolderListResponse,
  FormCreateRequest,
  FormCreateResponse,
  FormGetRequest,
  FormGetResponse,
  FormListResponse,
  RecordCreateRequest,
  RecordCreateResponse,
  RecordGetRequest,
  RecordGetResponse,
  RecordListRequest,
  RecordListResponse,
} from './types/client';

export default class Client extends BaseRESTClient implements ClientDefinition {
  static corev1 = 'apis/core.nrc.no/v1';

  constructor(address: string) {
    super(`${address}/${Client.corev1}`);
  }

  public Database = {
    create: (
      request: DatabaseCreateRequest,
    ): Promise<DatabaseCreateResponse> => {
      return this.do(request, '/databases', 'post', request.object, 200);
    },
    list: (request: {} | undefined): Promise<DatabaseListResponse> => {
      return this.do(request, '/databases', 'get', undefined, 200);
    },
  };

  public Folder = {
    create: (request: FolderCreateRequest): Promise<FolderCreateResponse> => {
      return this.do(request, '/folders', 'post', request.object, 200);
    },
    list: (request: {} | undefined): Promise<FolderListResponse> => {
      return this.do(request, '/folders', 'get', undefined, 200);
    },
  };

  public Form = {
    create: (request: FormCreateRequest): Promise<FormCreateResponse> => {
      return this.do(request, '/forms', 'post', request.object, 200);
    },
    list: (request: {} | undefined): Promise<FormListResponse> => {
      return this.do(request, '/forms', 'get', undefined, 200);
    },
    get: (request: FormGetRequest): Promise<FormGetResponse> => {
      return this.do(request, `/forms/${request.id}`, 'get', undefined, 200);
    },
  };

  public Record = {
    create: (request: RecordCreateRequest): Promise<RecordCreateResponse> => {
      return this.do(request, '/records', 'post', request.object, 200);
    },
    list: (request: RecordListRequest): Promise<RecordListResponse> => {
      const { databaseId, formId } = request;
      const url = `/records?databaseId=${databaseId}&formId=${formId}`;
      return this.do(request, url, 'get', undefined, 200);
    },
    get: (request: RecordGetRequest): Promise<RecordGetResponse> => {
      const { databaseId, formId, recordId } = request;
      const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
      return this.do(request, url, 'get', undefined, 200);
    },
  };
}
