import { AxiosInstance } from "axios";
export declare class Database {
    id: string;
    name: string;
}
export declare class DatabaseList {
    items: Database[];
}
export declare enum FieldKind {
    Text = "text",
    MultilineText = "multilineText",
    Reference = "reference",
    SubForm = "subform",
    Date = "date",
    Quantity = "quantity",
    SingleSelect = "singleSelect"
}
export declare class FieldType {
    text?: FieldTypeText;
    reference?: FieldTypeReference;
    subForm?: FieldTypeSubForm;
    multilineText?: FieldTypeMultilineText;
    date?: FieldTypeDate;
    quantity?: FieldTypeQuantity;
    singleSelect?: FieldTypeSingleSelect;
}
export declare class FieldTypeText {
}
export declare class FieldTypeMultilineText {
}
export declare class FieldTypeDate {
}
export declare class FieldTypeQuantity {
}
export declare class FieldTypeSingleSelect {
}
export declare class FieldTypeReference {
    databaseId: string;
    formId: string;
}
export declare class FieldTypeSubForm {
    id: string;
    name: string;
    code: string;
    fields: FieldDefinition[];
}
export declare class FieldDefinition {
    id: string;
    code: string;
    name: string;
    description: string;
    required: boolean;
    options: string[];
    key: boolean;
    fieldType: FieldType;
}
export declare class FormDefinition {
    id: string;
    code: string;
    databaseId: string;
    folderId: string;
    name: string;
    fields: FieldDefinition[];
}
export declare class FormDefinitionList {
    items: FormDefinition[];
}
export declare class Folder {
    id: string;
    databaseId: string;
    parentId: string;
    name: string;
}
export declare class FolderList {
    items: Folder[];
}
export declare class Record {
    id: string;
    databaseId: string;
    formId: string;
    parentId: string | undefined;
    values: {
        [key: string]: any;
    };
}
export declare class LocalRecord extends Record {
    isNew: boolean;
}
export declare type RecordList = {
    items: Record[];
};
export declare type Response<TRequest, TResponse> = {
    request: TRequest;
    response: TResponse | undefined;
    status: string;
    statusCode: number;
    success: boolean;
    error: any;
};
export declare type PartialObjectWrapper<T> = {
    object: Partial<T>;
};
export declare type DataOperation<TRequest, TResponse> = (request: TRequest) => Promise<TResponse>;
export declare type DatabaseCreateRequest = PartialObjectWrapper<Database>;
export declare type DatabaseCreateResponse = Response<DatabaseCreateRequest, Database>;
export interface DatabaseCreator {
    createDatabase: DataOperation<DatabaseCreateRequest, DatabaseCreateResponse>;
}
export declare type DatabaseListRequest = {} | undefined;
export declare type DatabaseListResponse = Response<DatabaseListRequest, DatabaseList>;
export interface DatabaseLister {
    listDatabases: DataOperation<DatabaseListRequest, DatabaseListResponse>;
}
export declare type FormListRequest = {} | undefined;
export declare type FormListResponse = Response<FormListRequest, FormDefinitionList>;
export declare type FormGetRequest = {
    id: string;
};
export declare type FormGetResponse = Response<FormGetRequest, FormDefinition>;
export interface FormGetter {
    getForm: DataOperation<FormGetRequest, FormGetResponse>;
}
export interface FormLister {
    listForms: DataOperation<FormListRequest, FormListResponse>;
}
export declare type FormCreateRequest = PartialObjectWrapper<FormDefinition>;
export declare type FormCreateResponse = Response<FormCreateRequest, FormDefinition>;
export interface FormCreator {
    createForm: DataOperation<FormCreateRequest, FormCreateResponse>;
}
export declare type RecordCreateRequest = PartialObjectWrapper<Record>;
export declare type RecordCreateResponse = Response<RecordCreateRequest, Record>;
export interface RecordCreator {
    createRecord: DataOperation<RecordCreateRequest, RecordCreateResponse>;
}
export declare type FolderListRequest = {} | undefined;
export declare type FolderListResponse = Response<FolderListRequest, FolderList>;
export interface FolderLister {
    listFolders: DataOperation<FolderListRequest, FolderListResponse>;
}
export declare type FolderCreateRequest = PartialObjectWrapper<Folder>;
export declare type FolderCreateResponse = Response<FolderCreateRequest, Folder>;
export interface FolderCreator {
    createFolder: DataOperation<FolderCreateRequest, FolderCreateResponse>;
}
export declare type RecordListRequest = {
    databaseId: string;
    formId: string;
};
export declare type RecordListResponse = Response<RecordListRequest, RecordList>;
export interface RecordLister {
    listRecords: DataOperation<RecordListRequest, RecordListResponse>;
}
export declare type RecordGetRequest = {
    databaseId: string;
    formId: string;
    recordId: string;
};
export declare type RecordGetResponse = Response<RecordGetRequest, Record>;
export interface RecordGetter {
    getRecord: DataOperation<RecordGetRequest, RecordGetResponse>;
}
export interface ClientDefinition extends DatabaseCreator, DatabaseLister, FormLister, FormGetter, FormCreator, RecordCreator, RecordLister, RecordGetter, FolderLister, FolderCreator {
    address: string;
    axiosInstance: AxiosInstance;
}
