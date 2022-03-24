import { FormDefinition, FormWithRecord, Record, RecordList } from '../model';

import { FormLookup } from './Form';
import { DataOperation, Response } from './utils';
import { RecordLookup } from './Record';

export type Recipient = Record;
export type RecipientList = RecordList;
export type RecipientListResponse = Response<FormLookup, RecipientList>;

export interface RecipientClientDefinition {
  create: DataOperation<
    FormWithRecord<Recipient>[],
    FormWithRecord<Recipient>[]
  >;
  list: DataOperation<FormLookup, FormWithRecord<Recipient>[][]>;
  get: DataOperation<RecordLookup, FormWithRecord<Recipient>[]>;
  getRecipientForms: () => Promise<FormDefinition[]>;
}
