import { AxiosInstance, Method } from "axios";
import { ClientDefinition, DatabaseCreateRequest, DatabaseCreateResponse, DatabaseListResponse, FolderCreateRequest, FolderCreateResponse, FolderListResponse, FormCreateRequest, FormCreateResponse, FormGetRequest, FormGetResponse, FormListResponse, RecordCreateRequest, RecordCreateResponse, RecordGetRequest, RecordGetResponse, RecordListRequest, RecordListResponse, Response } from "./types/types";
export default class client implements ClientDefinition {
    readonly address: string;
    readonly axiosInstance: AxiosInstance;
    constructor(address?: string, axiosInstance?: AxiosInstance);
    do<TRequest, TBody>(request: TRequest, url: string, method: Method, data: any, expectStatusCode: number): Promise<Response<TRequest, TBody>>;
    createDatabase(request: DatabaseCreateRequest): Promise<DatabaseCreateResponse>;
    createFolder(request: FolderCreateRequest): Promise<FolderCreateResponse>;
    createForm(request: FormCreateRequest): Promise<FormCreateResponse>;
    createRecord(request: RecordCreateRequest): Promise<RecordCreateResponse>;
    listDatabases(request: {} | undefined): Promise<DatabaseListResponse>;
    listFolders(request: {} | undefined): Promise<FolderListResponse>;
    listForms(request: {} | undefined): Promise<FormListResponse>;
    listRecords(request: RecordListRequest): Promise<RecordListResponse>;
    getForm(request: FormGetRequest): Promise<FormGetResponse>;
    getRecord(request: RecordGetRequest): Promise<RecordGetResponse>;
}
export declare const defaultClient: ClientDefinition;
