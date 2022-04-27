import {
  FieldDefinition,
  FieldValue,
  FormLookup,
  FormWithRecord,
  RecordLookup,
  Record,
  FormDefinition,
  FormType,
} from '../types';
import {
  Recipient,
  RecipientClientDefinition,
} from '../types/client/Recipient';
import { Validation } from '..';
import * as Tree from '../utils/tree';

import { RecordClient } from './Record';
import { FormClient } from './Form';
import * as console from "console";

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
    Validation.Recipient.validateRecipientHierarchy(recipient);

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

  public list = async ({
    formId,
    databaseId,
  }: FormLookup): Promise<FormWithRecord<Recipient>[][]> => {
    const response = await this.recordClient.list({
      formId,
      databaseId,
      subforms: false,
    });
    if (!response.response) {
      return response.error;
    }
    return Promise.all(
      response.response?.items.map((item) => {
        return this.get({ formId, databaseId, recordId: item.id });
      }),
    );
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

  public getRecipientForms = async (): Promise<FormDefinition[]> => {
    const formResponse = await this.formClient.list(undefined);

    if (!formResponse.response) {
      throw new Error(formResponse.error);
    }

    const recipeintForms = formResponse.response.items
      .filter((f) => f.formType === FormType.RecipientFormType)
      .map((f) => {
        const referenceKey = f.fields.find(
          (field: FieldDefinition) => field.key && field.fieldType.reference,
        );
        return {
          ...f,
          parentId: referenceKey?.fieldType.reference?.formId || '',
        };
      });

    const tree = Tree.createDataTree(recipeintForms, 'id', 'parentId');
    const leaves = Tree.getLeafNodes(tree);
    return leaves.map((l) => {
      const { parentId: _, childNodes: __, ...f } = l;
      return f;
    });
  };
}
