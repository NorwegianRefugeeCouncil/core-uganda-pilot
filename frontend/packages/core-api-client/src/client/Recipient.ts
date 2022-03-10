import {
  FieldDefinition,
  FieldValue,
  PopulatedForm,
  RecordLookup,
} from '../types';
import {
  Recipient,
  RecipientDefinition,
  RecipientList,
  RecipientResponse,
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

  get = (args: RecordLookup): Promise<RecipientResponse> => {
    return this.recordClient.get(args);
  };

  getAncestors = async ({
    recordId,
    formId,
    databaseId,
  }: RecordLookup): Promise<PopulatedForm<Recipient>[]> => {
    const formResponse = await this.formClient.get({ id: formId });

    if (!formResponse.response) {
      throw new Error(formResponse.error);
    }

    const reference = formResponse.response.fields.find(
      (field: FieldDefinition) => field.fieldType.reference,
    );

    const recipientGetResponse = await this.get({
      recordId,
      formId,
      databaseId,
    });

    if (!recipientGetResponse.response) {
      throw new Error(recipientGetResponse.error);
    }

    if (reference && reference.fieldType.reference) {
      const parentRecord = recipientGetResponse.response.values.find(
        (v: FieldValue) => v.fieldId === reference.id,
      );
      const parentRecordId = parentRecord?.value;

      if (!parentRecordId || !reference.fieldType.reference) {
        throw new Error('broken reference');
      }
      const result = await this.getAncestors({
        recordId: parentRecordId as string,
        formId: reference.fieldType.reference.formId,
        databaseId: reference.fieldType.reference.databaseId,
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
