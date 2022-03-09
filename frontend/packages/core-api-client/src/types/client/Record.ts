import { Record, RecordList } from '../model';

import { FormLookup } from './Form';
import { DataOperation, PartialObjectWrapper, Response } from './utils';

export type RecordCreateRequest = PartialObjectWrapper<Record>;
export type RecordCreateResponse = Response<RecordCreateRequest, Record>;

export type RecordListResponse = Response<FormLookup, RecordList>;

export type RecordLookup = FormLookup & { recordId: string };

export type RecordListRequest = { databaseId: string; formId: string };

export type RecordGetRequest = {
  databaseId: string;
  formId: string;
  recordId: string;
};
export type RecordGetResponse = Response<RecordGetRequest, Record>;

export interface RecordClientDefinition {
  create: DataOperation<RecordCreateRequest, RecordCreateResponse>;
  list: DataOperation<RecordListRequest, RecordListResponse>;
  get: DataOperation<RecordGetRequest, RecordGetResponse>;
}
