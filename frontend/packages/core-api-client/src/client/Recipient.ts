import {
  FieldDefinition,
  FieldValue,
  FormWithRecord,
  RecordLookup,
  Record,
} from '../types';
import {
  Recipient,
  RecipientClientDefinition,
  RecipientList,
} from '../types/client/Recipient';
import { Validation } from '..';

import { RecordClient } from './Record';
import { FormClient } from './Form';

export class RecipientClient implements RecipientClientDefinition {
  private recordClient: RecordClient;

  private formClient: FormClient;

  constructor(recordClient: RecordClient, formClient: FormClient) {
    this.recordClient = recordClient;
    this.formClient = formClient;
  }

  public create = (
    recipient: FormWithRecord<Recipient>[],
  ): Promise<FormWithRecord<Recipient>[]> => {
    // Throws on validation error
    Validation.Recipient.validateRecipientHierarch(recipient);

    return recipient.reduce<Promise<FormWithRecord<Recipient>[]>>(
      async (acc, { form, record }, idx) => {
        const resolvedAcc = await acc;

        // If the record is the root we just create it
        // If it is a child we first set the key value to the parent that was just created
        const keyFieldId = form.fields.find((f) => f.key)?.id;
        const parsedRecord: Record =
          idx === 0
            ? record
            : {
                ...record,
                values: record.values.map((v) => {
                  if (v.fieldId === keyFieldId)
                    return {
                      fieldId: v.fieldId,
                      value: resolvedAcc[idx - 1].record.id,
                    };
                  return v;
                }),
              };

        const createdRecord = await this.recordClient.createWithSubRecords({
          form,
          record: parsedRecord,
        });

        return [...resolvedAcc, { form, record: createdRecord }];
      },
      Promise.resolve([]),
    );
  };

  public list = (): Promise<RecipientList> => {
    return Promise.resolve({ items: [] });
  };

  public get = async ({
    recordId,
    formId,
    databaseId,
  }: RecordLookup): Promise<FormWithRecord<Recipient>[]> => {
    const formResponse = await this.formClient.get({ id: formId });

    if (!formResponse.response) {
      throw new Error(formResponse.error);
    }

    const referenceKey = formResponse.response.fields.find(
      (field: FieldDefinition) => field.key && field.fieldType.reference,
    );

    const recipientGetResponse = await this.recordClient.get({
      recordId,
      formId,
      databaseId,
    });

    if (!recipientGetResponse.response) {
      throw new Error(recipientGetResponse.error);
    }

    if (referenceKey) {
      const parentRecord = recipientGetResponse.response.values.find(
        (v: FieldValue) => v.fieldId === referenceKey.id,
      );
      const parentRecordId = parentRecord?.value;

      if (!parentRecordId || !referenceKey.fieldType.reference) {
        throw new Error('broken reference');
      }
      const result = await this.get({
        recordId: parentRecordId as string,
        formId: referenceKey.fieldType.reference.formId,
        databaseId: referenceKey.fieldType.reference.databaseId,
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
