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
} from './types';

export default class Client extends BaseRESTClient implements ClientDefinition {
  static corev1 = 'apis/core.nrc.no/v1';

  constructor(address: string) {
    super(`${address}/${Client.corev1}`);
  }

  createDatabase(request: DatabaseCreateRequest): Promise<DatabaseCreateResponse> {
    return this.do(request, '/databases', 'post', request.object, 200);
  }

  createFolder(request: FolderCreateRequest): Promise<FolderCreateResponse> {
    return this.do(request, '/folders', 'post', request.object, 200);
  }

  createForm(request: FormCreateRequest): Promise<FormCreateResponse> {
    return this.do(request, '/forms', 'post', request.object, 200);
  }

  createRecord(request: RecordCreateRequest): Promise<RecordCreateResponse> {
    const url = '/records';
    return this.do(request, url, 'post', request.object, 200);
  }

  listDatabases(request: {} | undefined): Promise<DatabaseListResponse> {
    return this.do(request, '/databases', 'get', undefined, 200);
  }

  listFolders(request: {} | undefined): Promise<FolderListResponse> {
    return this.do(request, '/folders', 'get', undefined, 200);
  }

  listForms(request: {} | undefined): Promise<FormListResponse> {
    return this.do(request, '/forms', 'get', undefined, 200);
  }

  listRecords(request: RecordListRequest): Promise<RecordListResponse> {
    const { databaseId, formId } = request;
    const url = `/records?databaseId=${databaseId}&formId=${formId}`;
    return this.do(request, url, 'get', undefined, 200);
  }

  getForm(request: FormGetRequest): Promise<FormGetResponse> {
    return this.do(request, `/forms/${request.id}`, 'get', undefined, 200);
  }

  getRecord(request: RecordGetRequest): Promise<RecordGetResponse> {
    const { databaseId, formId, recordId } = request;
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    return this.do(request, url, 'get', undefined, 200);
  }
}
