import { FormWithRecord, Record, RecordList } from '../model';

import { FormLookup } from './Form';
import { DataOperation, Response } from './utils';
import { RecordLookup } from './Record';

export type Recipient = Record;
export type RecipientDefinition = Omit<Recipient, 'id'>;
export type RecipientList = RecordList;

export interface RecipientClientDefinition {
  create: DataOperation<RecipientDefinition, Recipient>;
  list: DataOperation<FormLookup, RecipientList>;
  get: DataOperation<RecordLookup, FormWithRecord<Recipient>[]>;
}
