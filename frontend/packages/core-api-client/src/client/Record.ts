import { getFieldKind } from '..';
import {
  Record,
  FieldKind,
  RecordCreateRequest,
  RecordCreateResponse,
  RecordGetRequest,
  RecordGetResponse,
  RecordListRequest,
  RecordListResponse,
  RecordList,
} from '../types';

import { BaseRESTClient } from './BaseRESTClient';
import { FormClient } from './Form';

export class RecordClient {
  restClient: BaseRESTClient;

  formClient: FormClient;

  constructor(restClient: BaseRESTClient, formClient: FormClient) {
    this.restClient = restClient;
    this.formClient = formClient;
  }

  create = (request: RecordCreateRequest): Promise<RecordCreateResponse> => {
    return this.restClient.do(request, '/records', 'post', request.object, 200);
  };

  list = async (request: RecordListRequest): Promise<RecordListResponse> => {
    const { databaseId, formId } = request;
    const url = `/records?databaseId=${databaseId}&formId=${formId}`;
    const response = await this.restClient.do<RecordListRequest, RecordList>(
      request,
      url,
      'get',
      undefined,
      200,
    );

    if (!response.response) return response;

    const recordList = await Promise.all(
      response.response.items.map(async (record) => this.massageRecord(record)),
    );

    return {
      ...response,
      response: {
        ...response.response,
        items: recordList,
      },
    };
  };

  get = async (request: RecordGetRequest): Promise<RecordGetResponse> => {
    const { databaseId, formId, recordId } = request;
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    const response = await this.restClient.do<RecordGetRequest, Record>(
      request,
      url,
      'get',
      undefined,
      200,
    );

    if (!response.response) return response;

    const record = await this.massageRecord(response.response);

    return {
      ...response,
      response: record,
    };
  };

  /**
   * Takes a record and populates it with additional information
   * Fetches subrecord values and adds them to the record's values
   * TODO: Change boolean values from string to booleans
   * TODO: Fetch reference field values to have a human readable name
   */
  private massageRecord = async (record: Record): Promise<Record> => {
    // Fetching a subrecords form returns the owner form causing an infinite loop so we skip over subrecords
    if (record.ownerId) return record;

    // Fetch the record's form to get a complete list of fields
    const formResponse = await this.formClient.get({ id: record.formId });

    if (!formResponse.response) return record;

    // Iterate over the forms fields
    // if it's a subform field, fetch the subform records and populate with their values
    // otherwise, just add the field value to the record
    const values = await Promise.all(
      formResponse.response?.fields.map(async (field) => {
        if (getFieldKind(field.fieldType) === FieldKind.SubForm) {
          const subrecordList = await this.list({
            databaseId: record.databaseId,
            formId: field.id,
          });
          return {
            fieldId: field.id,
            value:
              subrecordList.response?.items.map(
                (subrecord) => subrecord.values,
              ) ?? [],
          };
        }

        const existingValue = record.values.find(
          (value) => value.fieldId === field.id,
        );
        if (!existingValue) throw new Error(); // Should never happen, unless form and record don't match
        return existingValue;
      }) ?? [],
    );

    return {
      ...record,
      values,
    };
  };
}
