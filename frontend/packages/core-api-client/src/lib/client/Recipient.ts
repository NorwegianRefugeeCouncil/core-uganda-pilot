import { BaseRESTClient } from './BaseRESTClient';
import {Response} from "../types";
import {
  Recipient,
  RecipientDefinition, RecipientGetRequest,
  RecipientList, RecipientListRequest
} from "../types/client/Recipient";

export class RecordClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = (recipient: RecipientDefinition): Promise<Response<RecipientDefinition, Recipient>> => {
    return this.restClient.post( '/records', recipient);
  };

  list = (request: RecipientListRequest): Promise<Response<undefined, RecipientList>> => {
    const { databaseId, formId } = request;
    const url = `/records?databaseId=${databaseId}&formId=${formId}`;
    return this.restClient.get(url);
  };

  get = async (request: RecipientGetRequest): Promise<Response<undefined, Recipient>> => {
    const { databaseId, formId, recordId } = request;
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    return this.restClient.get(url);
  };
}
