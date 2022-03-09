import {
  ErrorResponse,
  FormLookup,
  RecordGetRequest,
  RecordLookup,
} from '../types';
import {
  Recipient,
  RecipientDefinition,
  RecipientList, RecipientResponse,
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

  get = (args: RecordLookup): Promise<RecipientResponse> => {
    return this.recordClient.get(args);
  };
}
