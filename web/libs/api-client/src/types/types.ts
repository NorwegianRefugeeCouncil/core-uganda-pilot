import {AxiosInstance} from "axios";

export class Database {
    public id: string = ""
    public name: string = ""
}

export class DatabaseList {
    public items: Database[] = []
}

export enum FieldKind {
    Text = "text",
    MultilineText = "multilineText",
    Reference = "reference",
    SubForm = "subform",
    Date = "date",
    Quantity = "quantity",
    SingleSelect = "singleSelect",
}

export class FieldType {
    public text?: FieldTypeText
    public reference?: FieldTypeReference
    public subForm?: FieldTypeSubForm
    public multilineText?: FieldTypeMultilineText
    public date?: FieldTypeDate
    public quantity?: FieldTypeQuantity
    public singleSelect?: FieldTypeSingleSelect
}

export class FieldTypeText {
}

export class FieldTypeMultilineText {
}

export class FieldTypeDate {
}

export class FieldTypeQuantity {
}

export class FieldTypeSingleSelect {
}

export class FieldTypeReference {
    public databaseId: string = ""
    public formId: string = ""
}

export class FieldTypeSubForm {
    public id: string = ""
    public name: string = ""
    public code: string = ""
    public fields: FieldDefinition[] = []
}

export class FieldDefinition {
    public id: string = ""
    public code: string = ""
    public name: string = ""
    public description: string = ""
    public required: boolean = false
    public options: string[] = []
    public key: boolean = false
    public fieldType: FieldType = new FieldType()
}

export class FormDefinition {
    public id: string = ""
    public code: string = ""
    public databaseId: string = ""
    public folderId: string = ""
    public name: string = ""
    public fields: FieldDefinition[] = []
}

export class FormDefinitionList {
    public items: FormDefinition[] = []
}

export class Folder {
    public id: string = ""
    public databaseId: string = ""
    public parentId: string = ""
    public name: string = ""
}

export class FolderList {
    public items: Folder[] = []
}

export class Record {
    public id: string = ""
    public databaseId: string = ""
    public formId: string = ""
    public parentId: string | undefined = undefined
    public values: { [key: string]: any } = {}
}

export class LocalRecord extends Record {
    public isNew: boolean = false
}

export type RecordList = { items: Record[] }


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

export type FormGetRequest = { id: string }
export type FormGetResponse = Response<FormGetRequest, FormDefinition>

export interface FormGetter {
    getForm: DataOperation<FormGetRequest, FormGetResponse>
}

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

export type RecordGetRequest = { databaseId: string, formId: string, recordId: string }
export type RecordGetResponse = Response<RecordGetRequest, Record>

export interface RecordGetter {
    getRecord: DataOperation<RecordGetRequest, RecordGetResponse>
}


export interface ClientDefinition
    extends DatabaseCreator,
        DatabaseLister,
        FormLister,
        FormGetter,
        FormCreator,
        RecordCreator,
        RecordLister,
        RecordGetter,
        FolderLister,
        FolderCreator {
    address: string
    axiosInstance: AxiosInstance
}
