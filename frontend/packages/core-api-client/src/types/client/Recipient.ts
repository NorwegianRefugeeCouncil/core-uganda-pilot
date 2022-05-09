import { FormDefinition, FormWithRecord, Record, RecordList } from '../model';

import { FormLookup } from './Form';
import { DataOperation } from './utils';
import { RecordLookup } from './Record';

export type Recipient = Record;
export type RecipientList = RecordList;

export interface RecipientClientDefinition {
  create: DataOperation<
    FormWithRecord<Recipient>[],
    FormWithRecord<Recipient>[]
  >;
  list: DataOperation<FormLookup, RecipientList>;
  get: DataOperation<RecordLookup, FormWithRecord<Recipient>[]>;
  getRecipientForms: () => Promise<FormDefinition[]>;
}
