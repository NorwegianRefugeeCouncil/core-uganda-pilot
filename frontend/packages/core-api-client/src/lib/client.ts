import axios, { AxiosInstance, Method } from 'axios';

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
  Response,
} from './types';
import { clientResponse } from './utils/responses';

export default class Client implements ClientDefinition {
  private readonly axiosInstance: AxiosInstance;

  private readonly corev1 = 'apis/core.nrc.no/v1';

  constructor(address: string) {
    this.axiosInstance = axios.create({
      baseURL: `${address}/${this.corev1}`,
    });
  }

  public setAuth = (token: string): void => {
    this.axiosInstance.interceptors.request.use((value: any) => {
      const result = { ...value };
      if (!token) {
        return value;
      }
      return {
        ...result,
        headers: {
          ...result.headers,
          Authorization: `Bearer ${token}`,
        },
      };
    });
  };

  do<TRequest, TBody>(
    request: TRequest,
    url: string,
    method: Method,
    data: any,
    expectStatusCode: number,
  ): Promise<Response<TRequest, TBody>> {
    const headers: { [key: string]: string } = {
      Accept: 'application/json',
    };
    return this.axiosInstance
      .request<TBody>({
        responseType: 'json',
        method,
        url,
        data,
        headers,
        withCredentials: true,
      })
      .then((value) => {
        return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
      });
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
