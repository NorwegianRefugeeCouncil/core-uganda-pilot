import { Record, RecordList } from '../model';

import { DataOperation } from './utils';
import { RecordLookup, FormLookup } from './Record';

export type Recipient = Record;
export type RecipientDefinition = Omit<Recipient, 'id'>;
export type RecipientList = RecordList;

export interface RecipientClientDefinition {
  create: DataOperation<RecipientDefinition, Recipient>;
  list: DataOperation<FormLookup, RecipientList>;
  get: DataOperation<RecordLookup, Recipient>;
}
