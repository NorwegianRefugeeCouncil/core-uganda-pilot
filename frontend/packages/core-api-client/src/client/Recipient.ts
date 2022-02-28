import { FormLookup, RecordLookup } from '../types';
import {
  Recipient,
  RecipientDefinition,
  RecipientList,
} from '../types/client/Recipient';

import { RecordClient } from './Record';

export class RecipientClient {
  recordClient: RecordClient;

  constructor(recordClient: RecordClient) {
    this.recordClient = recordClient;
  }

  create = (recipient: RecipientDefinition): Promise<Recipient> => {
    return Promise.resolve({
      id: 'id',
      ...recipient,
    });
  };

  list = (args: FormLookup): Promise<RecipientList> => {
    return Promise.resolve({ items: [] });
  };

  get = async (args: RecordLookup): Promise<Recipient> => {
    return Promise.resolve({
      id: 'id',
      formId: 'formId',
      databaseId: 'databaseId',
      ownerId: undefined,
      values: [],
    });
  };
}
