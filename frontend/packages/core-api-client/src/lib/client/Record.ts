import { BaseRESTClient } from './BaseRESTClient';
import {
  Record,
  RecordLookup,
  RecordList,
  FormLookup, RecordDefinition
} from "../types";

export class RecordClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = async (record: RecordDefinition): Promise<Record> => {
    const apiResponse = await this.restClient.post('/records', record);
    if (apiResponse.success) {
      return apiResponse.response as Record
    }
    return apiResponse.error
  };

  list = async ({ databaseId, formId }: FormLookup): Promise<RecordList> => {
    const url = `/records?databaseId=${databaseId}&formId=${formId}`;
    const apiResponse = await this.restClient.get(url);
    if (apiResponse.success) {
      return apiResponse.response as RecordList
    }
    return apiResponse.error
  };

  get = async ({ databaseId, formId, recordId }: RecordLookup): Promise<Record> => {
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    const apiResponse = await this.restClient.get(url);
    if (apiResponse.success) {
      return apiResponse.response as Record
    }
    return apiResponse.error
  };
}
