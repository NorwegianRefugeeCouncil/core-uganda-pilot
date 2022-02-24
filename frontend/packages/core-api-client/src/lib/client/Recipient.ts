import {FormLookup, RecordLookup, Response} from "../types";
import {
  Recipient,
  RecipientDefinition,
  RecipientList,
} from "../types/client/Recipient";
import { RecordClient } from "./Record";

export class RecipientClient {
  recordClient: RecordClient;

  constructor(recordClient: RecordClient) {
    this.recordClient = recordClient;
  }

  create = (recipient: RecipientDefinition): Promise<Recipient> => {
    return this.recordClient.create(recipient);
  };

  list = (args: FormLookup): Promise<RecipientList> => {
    return this.recordClient.list(args);
  };

  get = async (args: RecordLookup): Promise<Recipient> => {
    return this.recordClient.get(args);
  };
}
