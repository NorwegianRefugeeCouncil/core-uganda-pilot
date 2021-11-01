import axios, {AxiosResponse, Method} from "axios";
import {
    Database,
    DatabaseList,
    Folder,
    FolderList,
    FormDefinition,
    FormDefinitionList,
    Record,
    RecordList
} from "./types/types";
import {useAuth} from "../../../admin/src/app/hooks";

export type Response<TRequest, TResponse> = {
    request: TRequest
    response: TResponse | undefined
    status: string
    statusCode: number
    success: boolean
    error: any
}

export type PartialObjectWrapper<T> = { object: Partial<T> }
export type DataOperation<TRequest, TResponse> = (request: TRequest) => Promise<TResponse>

export type DatabaseCreateRequest = PartialObjectWrapper<Database>
export type DatabaseCreateResponse = Response<DatabaseCreateRequest, Database>

export interface DatabaseCreator {
    createDatabase: DataOperation<DatabaseCreateRequest, DatabaseCreateResponse>
}

export type DatabaseListRequest = {} | undefined
export type DatabaseListResponse = Response<DatabaseListRequest, DatabaseList>

export interface DatabaseLister {
    listDatabases: DataOperation<DatabaseListRequest, DatabaseListResponse>
}

export type FormListRequest = {} | undefined
export type FormListResponse = Response<FormListRequest, FormDefinitionList>

export interface FormLister {
    listForms: DataOperation<FormListRequest, FormListResponse>
}

export type FormCreateRequest = PartialObjectWrapper<FormDefinition>
export type FormCreateResponse = Response<FormCreateRequest, FormDefinition>

export interface FormCreator {
    createForm: DataOperation<FormCreateRequest, FormCreateResponse>
}

export type RecordCreateRequest = PartialObjectWrapper<Record>
export type RecordCreateResponse = Response<RecordCreateRequest, Record>

export interface RecordCreator {
    createRecord: DataOperation<RecordCreateRequest, RecordCreateResponse>
}

export type FolderListRequest = {} | undefined
export type FolderListResponse = Response<FolderListRequest, FolderList>

export interface FolderLister {
    listFolders: DataOperation<FolderListRequest, FolderListResponse>
}

export type FolderCreateRequest = PartialObjectWrapper<Folder>
export type FolderCreateResponse = Response<FolderCreateRequest, Folder>

export interface FolderCreator {
    createFolder: DataOperation<FolderCreateRequest, FolderCreateResponse>
}

export type RecordListRequest = { databaseId: string, formId: string }
export type RecordListResponse = Response<RecordListRequest, RecordList>

export interface RecordLister {
    listRecords: DataOperation<RecordListRequest, RecordListResponse>
}

export interface RecordLister {
    listRecords: DataOperation<RecordListRequest, RecordListResponse>
}


export interface Client
    extends DatabaseCreator,
        DatabaseLister,
        FormLister,
        FormCreator,
        RecordCreator,
        RecordLister,
        FolderLister,
        FolderCreator {
}

function errorResponse<TRequest, TBody>(request: TRequest, r: AxiosResponse<TBody>): Response<TRequest, TBody> {
    return {
        request: request,
        response: undefined,
        status: r.request,
        statusCode: r.status,
        error: r.data as any,
        success: false,
    };
}

function successResponse<TRequest, TBody>(request: TRequest, r: AxiosResponse<TBody>): Response<TRequest, TBody> {
    return {
        request: request,
        response: r.data as TBody,
        status: r.statusText,
        statusCode: r.status,
        error: undefined,
        success: true,
    };
}

function clientResponse<TRequest, TBody>(r: AxiosResponse<TBody>, request: TRequest, expectedStatusCode: number): Response<TRequest, TBody> {
    return r.status !== expectedStatusCode
        ? errorResponse<TRequest, TBody>(request, r)
        : successResponse<TRequest, TBody>(request, r)
}

export class client implements Client {
    public address = "http://localhost:9000"

    do<TRequest, TBody>(request: TRequest, url: string, method: Method, data: any, expectStatusCode: number): Promise<Response<TRequest, TBody>> {
        return axios.request<TBody>({
            method,
            url,
            data
        }).then(value => {
            return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
        }).catch((err) => {
            return {
                request: request,
                response: undefined,
                status: "500 Internal Server Error",
                statusCode: 500,
                error: err.message,
                success: false,
            }
        })
    }

    createDatabase(request: DatabaseCreateRequest): Promise<DatabaseCreateResponse> {
        return this.do(request, `${this.address}/databases`, "post", request.object, 200)
    }

    createFolder(request: FolderCreateRequest): Promise<FolderCreateResponse> {
        return this.do(request, `${this.address}/folders`, "post", request.object, 200)
    }

    createForm(request: FormCreateRequest): Promise<FormCreateResponse> {
        return this.do(request, `${this.address}/forms`, "post", request.object, 200)
    }

    createRecord(request: RecordCreateRequest): Promise<RecordCreateResponse> {
        const url = `${this.address}/records`
        return this.do(request, url, "post", request.object, 200)
    }

    listDatabases(request: {} | undefined): Promise<DatabaseListResponse> {
        return this.do(request, `${this.address}/databases`, "get", undefined, 200)
    }

    listFolders(request: {} | undefined): Promise<FolderListResponse> {
        return this.do(request, `${this.address}/folders`, "get", undefined, 200)
    }

    listForms(request: {} | undefined): Promise<FormListResponse> {
        return this.do(request, `${this.address}/forms`, "get", undefined, 200)
    }

    listRecords(request: RecordListRequest): Promise<RecordListResponse> {
        const {databaseId, formId} = request
        const url = `${this.address}/records?databaseId=${databaseId}&formId=${formId}`
        return this.do(request, url, "get", undefined, 200)
    }

}

export const defaultClient: Client = new client()

