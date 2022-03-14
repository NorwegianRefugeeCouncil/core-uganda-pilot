import {
  FieldDefinition,
  FieldValue,
  FormWithRecord,
  RecordLookup,
} from '../types';
import {
  Recipient,
  RecipientDefinition,
  RecipientList,
} from '../types/client/Recipient';

import { RecordClient } from './Record';
import { FormClient } from './Form';

export class RecipientClient {
  recordClient: RecordClient;

  formClient: FormClient;

  constructor(recordClient: RecordClient, formClient: FormClient) {
    this.recordClient = recordClient;
    this.formClient = formClient;
  }

  create = (recipient: RecipientDefinition): Promise<Recipient> => {
    return Promise.resolve({
      id: 'id',
      ...recipient,
    });
  };

  list = (): Promise<RecipientList> => {
    return Promise.resolve({ items: [] });
  };

  get = async ({
    recordId,
    formId,
    databaseId,
  }: RecordLookup): Promise<FormWithRecord<Recipient>[]> => {
    const formResponse = await this.formClient.get({ id: formId });

    if (!formResponse.response) {
      throw new Error(formResponse.error);
    }

    const key = formResponse.response.fields.find(
      (field: FieldDefinition) => field.key,
    );

    const recipientGetResponse = await this.recordClient.get({
      recordId,
      formId,
      databaseId,
    });

    if (!recipientGetResponse.response) {
      throw new Error(recipientGetResponse.error);
    }

    if (key && key.fieldType.reference) {
      const parentRecord = recipientGetResponse.response.values.find(
        (v: FieldValue) => v.fieldId === key.id,
      );
      const parentRecordId = parentRecord?.value;

      if (!parentRecordId || !key.fieldType.reference) {
        throw new Error('broken reference');
      }
      const result = await this.get({
        recordId: parentRecordId as string,
        formId: key.fieldType.reference.formId,
        databaseId: key.fieldType.reference.databaseId,
      });
      return [
        ...result,
        {
          form: formResponse.response,
          record: recipientGetResponse.response,
        },
      ];
    }
    return [
      {
        form: formResponse.response,
        record: recipientGetResponse.response,
      },
    ];
  };
}
