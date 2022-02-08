export type Database = {
  id: string;
  name: string;
};

export type DatabaseList = {
  items: Database[];
};

export enum FieldKind {
  Text = 'text',
  MultilineText = 'multilineText',
  Reference = 'reference',
  SubForm = 'subform',
  Date = 'date',
  Quantity = 'quantity',
  SingleSelect = 'singleSelect',
  MultiSelect = 'multiSelect',
  Week = 'week',
  Month = 'month',
  Checkbox = 'checkbox',
}

export type FieldType = {
  text?: FieldTypeText;
  reference?: FieldTypeReference;
  subForm?: FieldTypeSubForm;
  multilineText?: FieldTypeMultilineText;
  date?: FieldTypeDate;
  month?: FieldTypeMonth;
  week?: FieldTypeWeek;
  quantity?: FieldTypeQuantity;
  singleSelect?: FieldTypeSingleSelect;
  multiSelect?: FieldTypeMultiSelect;
  checkbox?: FieldTypeCheckbox;
};

export type FieldTypeText = {};

export type FieldTypeMultilineText = {};

export type FieldTypeDate = {};

export type FieldTypeMonth = {};

export type FieldTypeQuantity = {};

export type SelectOption = {
  name: string;
  id: string;
};

export type FieldTypeSingleSelect = {
  options: SelectOption[];
};

export type FieldTypeMultiSelect = {
  options: SelectOption[];
};

export class FieldTypeWeek {}

export type FieldTypeCheckbox = {};

export type FieldTypeReference = {
  databaseId: string;
  formId: string;
};

export type FieldTypeSubForm = {
  fields: FieldDefinition[];
};

export type FieldDefinition = {
  id: string;
  code: string;
  name: string;
  description: string;
  required: boolean;
  key: boolean;
  fieldType: FieldType;
};

export enum FormType {
  DefaultFormType = 'default',
  RecipientFormType = 'recipient',
}

export type FormDefinition = {
  id: string;
  code: string;
  databaseId: string;
  folderId: string;
  name: string;
  type: FormType;
  fields: FieldDefinition[];
};

export type FormDefinitionList = {
  items: FormDefinition[];
};

export type Folder = {
  id: string;
  databaseId: string;
  parentId: string;
  name: string;
};

export type FolderList = {
  items: Folder[];
};

export type FieldValue = {
  fieldId: string;
  value: string | string[] | null;
};

export type Record = {
  id: string;
  databaseId: string;
  formId: string;
  ownerId: string | undefined;
  values: FieldValue[];
};

export type LocalRecord = Record & {
  isNew: boolean;
};

export type RecordList = { items: Record[] };

export type RequestOptions = {
  headers: { [key: string]: string };
  silentRedirect?: boolean;
};

export type Response<TRequest, TResponse> = {
  request: TRequest;
  response: TResponse | undefined;
  status: string;
  statusCode: number;
  success: boolean;
  error: any;
};

export type PartialObjectWrapper<T> = { object: Partial<T> };
export type DataOperation<TRequest, TResponse> = (
  request: TRequest,
) => Promise<TResponse>;

export type DatabaseCreateRequest = PartialObjectWrapper<Database>;
export type DatabaseCreateResponse = Response<DatabaseCreateRequest, Database>;

export interface DatabaseCreator {
  createDatabase: DataOperation<DatabaseCreateRequest, DatabaseCreateResponse>;
}

export type DatabaseListRequest = {} | undefined;
export type DatabaseListResponse = Response<DatabaseListRequest, DatabaseList>;

export interface DatabaseLister {
  listDatabases: DataOperation<DatabaseListRequest, DatabaseListResponse>;
}

export type FormListRequest = {} | undefined;
export type FormListResponse = Response<FormListRequest, FormDefinitionList>;

export type FormGetRequest = { id: string };
export type FormGetResponse = Response<FormGetRequest, FormDefinition>;

export interface FormGetter {
  getForm: DataOperation<FormGetRequest, FormGetResponse>;
}

export interface FormLister {
  listForms: DataOperation<FormListRequest, FormListResponse>;
}

export type FormCreateRequest = PartialObjectWrapper<FormDefinition>;
export type FormCreateResponse = Response<FormCreateRequest, FormDefinition>;

export interface FormCreator {
  createForm: DataOperation<FormCreateRequest, FormCreateResponse>;
}

export type RecordCreateRequest = PartialObjectWrapper<Record>;
export type RecordCreateResponse = Response<RecordCreateRequest, Record>;

export interface RecordCreator {
  createRecord: DataOperation<RecordCreateRequest, RecordCreateResponse>;
}

export type FolderListRequest = {} | undefined;
export type FolderListResponse = Response<FolderListRequest, FolderList>;

export interface FolderLister {
  listFolders: DataOperation<FolderListRequest, FolderListResponse>;
}

export type FolderCreateRequest = PartialObjectWrapper<Folder>;
export type FolderCreateResponse = Response<FolderCreateRequest, Folder>;

export interface FolderCreator {
  createFolder: DataOperation<FolderCreateRequest, FolderCreateResponse>;
}

export type RecordListRequest = { databaseId: string; formId: string };
export type RecordListResponse = Response<RecordListRequest, RecordList>;

export interface RecordLister {
  listRecords: DataOperation<RecordListRequest, RecordListResponse>;
}

export type RecordGetRequest = {
  databaseId: string;
  formId: string;
  recordId: string;
};
export type RecordGetResponse = Response<RecordGetRequest, Record>;

export interface RecordGetter {
  getRecord: DataOperation<RecordGetRequest, RecordGetResponse>;
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
    FolderCreator {}
