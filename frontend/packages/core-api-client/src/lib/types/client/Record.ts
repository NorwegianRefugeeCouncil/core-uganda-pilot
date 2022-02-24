import {Record, RecordDefinition, RecordList} from '../model';

import { DataOperation, PartialObjectWrapper, Response } from './utils';

export type RecordCreateRequest = PartialObjectWrapper<Record>;
export type RecordCreateResponse = Response<RecordCreateRequest, Record>;

export type FormLookup = { databaseId: string; formId: string };
export type RecordListResponse = Response<FormLookup, RecordList>;

export type RecordLookup = FormLookup & { recordId: string };

export interface RecordClientDefinition {
  create: DataOperation<RecordDefinition, Record>;
  list: DataOperation<FormLookup , RecordList>;
  get: DataOperation<RecordLookup, Record>;
}
