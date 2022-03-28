import { getFieldKind } from '..';
import {
  Record,
  FieldKind,
  RecordCreateRequest,
  RecordCreateResponse,
  RecordGetResponse,
  RecordListRequest,
  RecordListResponse,
  RecordList,
  RecordLookup,
  FormDefinition,
  FieldValue,
  RecordClientDefinition,
  FormWithRecord,
} from '../types';

import { BaseRESTClient } from './BaseRESTClient';
import { FormClient } from './Form';

export class RecordClient implements RecordClientDefinition {
  private restClient: BaseRESTClient;

  private formClient: FormClient;

  static buildDefaultRecord = (form: FormDefinition): Record => ({
    id: '', // We are creating records so there shouldn't be an id
    databaseId: form.databaseId,
    formId: form.id,
    ownerId: undefined,
    values: form.fields.map((field) => {
      const fieldType = getFieldKind(field.fieldType);
      switch (fieldType) {
        case FieldKind.Text:
        case FieldKind.MultilineText:
        case FieldKind.Quantity:
          return {
            fieldId: field.id,
            value: '',
          };
        case FieldKind.Reference:
        case FieldKind.Date:
        case FieldKind.Month:
        case FieldKind.Week:
        case FieldKind.SingleSelect:
          return {
            fieldId: field.id,
            value: null,
          };
        case FieldKind.MultiSelect:
          return {
            fieldId: field.id,
            value: [],
          };
        case FieldKind.Checkbox:
          return {
            fieldId: field.id,
            value: 'false',
          };
        case FieldKind.SubForm:
          return {
            fieldId: field.id,
            value: [],
          };
        default:
          return {
            fieldId: field.id,
            value: '',
          };
      }
    }),
  });

  constructor(restClient: BaseRESTClient, formClient: FormClient) {
    this.restClient = restClient;
    this.formClient = formClient;
  }

  public create = (
    request: RecordCreateRequest,
  ): Promise<RecordCreateResponse> => {
    return this.restClient.do(request, '/records', 'post', request.object, 200);
  };

  // The backend does not create sub records, so we need to do it ourselves
  // This creates the record, then creates the sub records, then fetches the record again with it's subrecords
  public createWithSubRecords = async ({
    form,
    record,
  }: FormWithRecord<Record>): Promise<Record> => {
    const subFormFields = form.fields.filter(
      (f) => getFieldKind(f.fieldType) === FieldKind.SubForm,
    );

    const recordWithoutSubRecords = {
      ...record,
      values: record.values.filter(
        (v) => !subFormFields.some((sf) => sf.id === v.fieldId),
      ),
    };

    const recordResponse = await this.create({
      object: recordWithoutSubRecords,
    });

    if (recordResponse.error || !recordResponse.response)
      throw new Error(recordResponse.error);

    await Promise.all(
      subFormFields.map(async (f) => {
        const values = (record.values.find((v) => v.fieldId === f.id)?.value ??
          []) as FieldValue[][];
        const result = await Promise.all(
          values.map(async (v) => {
            const subRecord: Omit<Record, 'id'> = {
              ...record,
              formId: f.id,
              ownerId: recordResponse.response?.id,
              values: v,
            };
            const subRecordResponse = await this.create({ object: subRecord });
            if (subRecordResponse.error || !subRecordResponse.response)
              throw new Error(subRecordResponse.error);
            return subRecordResponse.response;
          }),
        );
        return [...result];
      }),
    );

    const getResponse = await this.get({
      recordId: recordResponse.response?.id,
      databaseId: recordResponse.response?.databaseId,
      formId: recordResponse.response?.formId,
    });

    if (getResponse.error || !getResponse.response)
      throw new Error(getResponse.error);

    return getResponse.response;
  };

  public list = async (
    request: RecordListRequest,
  ): Promise<RecordListResponse> => {
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

  public get = async (request: RecordLookup): Promise<RecordGetResponse> => {
    const { databaseId, formId, recordId } = request;
    const url = `/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
    const response = await this.restClient.do<RecordLookup, Record>(
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

  public buildDefaultRecord = (form: FormDefinition): Record =>
    RecordClient.buildDefaultRecord(form);

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
