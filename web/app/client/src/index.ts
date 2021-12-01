import axios, {AxiosInstance, Method} from "axios";
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
    Response
} from "./types/types";
import {clientResponse} from "./utils/responses";

export default class Client implements ClientDefinition {
    private corev1 = "apis/core.nrc.no/v1"

    constructor(
        public readonly address = 'https://core.dev:8443',
        public readonly axiosInstance: AxiosInstance = axios.create()) {
    }

    do<TRequest, TBody>(request: TRequest, url: string, method: Method, data: any, expectStatusCode: number): Promise<Response<TRequest, TBody>> {
        let headers: { [key: string]: string } = {
            "Accept": "application/json",
        }
        return this.axiosInstance.request<TBody>({
            responseType: "json",
            method,
            url,
            data,
            headers,
            withCredentials: true,
        }).then(value => {
            return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
        })
    }

    createDatabase(request: DatabaseCreateRequest): Promise<DatabaseCreateResponse> {
        return this.do(request, `${this.address}/${this.corev1}/databases`, "post", request.object, 200)
    }

    createFolder(request: FolderCreateRequest): Promise<FolderCreateResponse> {
        return this.do(request, `${this.address}/${this.corev1}/folders`, "post", request.object, 200)
    }

    createForm(request: FormCreateRequest): Promise<FormCreateResponse> {
        return this.do(request, `${this.address}/${this.corev1}/forms`, "post", request.object, 200)
    }

    createRecord(request: RecordCreateRequest): Promise<RecordCreateResponse> {
        const url = `${this.address}/${this.corev1}/records`
        return this.do(request, url, "post", request.object, 200)
    }

    listDatabases(request: {} | undefined): Promise<DatabaseListResponse> {
        return this.do(request, `${this.address}/${this.corev1}/databases`, "get", undefined, 200)
    }

    listFolders(request: {} | undefined): Promise<FolderListResponse> {
        return this.do(request, `${this.address}/${this.corev1}/folders`, "get", undefined, 200)
    }

    listForms(request: {} | undefined): Promise<FormListResponse> {
        return this.do(request, `${this.address}/${this.corev1}/forms`, "get", undefined, 200)
    }

    listRecords(request: RecordListRequest): Promise<RecordListResponse> {
        const {databaseId, formId} = request
        const url = `${this.address}/${this.corev1}/records?databaseId=${databaseId}&formId=${formId}`
        return this.do(request, url, "get", undefined, 200)
    }

    getForm(request: FormGetRequest): Promise<FormGetResponse> {
        return this.do(request, `${this.address}/${this.corev1}/forms/${request.id}`, "get", undefined, 200)
    }

    getRecord(request: RecordGetRequest): Promise<RecordGetResponse> {
        const {databaseId, formId, recordId} = request
        const url = `${this.address}/${this.corev1}/records/${recordId}?databaseId=${databaseId}&formId=${formId}`
        return this.do(request, url, "get", undefined, 200)
    }

}

export * from './types/types'
