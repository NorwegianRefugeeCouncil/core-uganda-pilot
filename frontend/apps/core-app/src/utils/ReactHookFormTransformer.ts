import {
  FieldKind,
  FieldValue,
  FormWithRecord,
  getFieldKind,
} from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FieldValue as RHFieldValue } from 'react-hook-form';

export const toReactHookForm = (
  data: FormWithRecord<Recipient>[],
): RHFieldValue<any> =>
  data.reduce((acc, { form, record }) => {
    return {
      ...acc,
      [form.id]: record.values.reduce((innerAcc, fieldValue) => {
        return {
          ...innerAcc,
          [fieldValue.fieldId]: fieldValue.value,
        };
      }, {}),
    };
  }, {});

export const fromReactHookForm = (
  originalData: FormWithRecord<Recipient>[],
  value: RHFieldValue<any>,
): FormWithRecord<Recipient>[] =>
  originalData.map((datum) => ({
    form: datum.form,
    record: {
      ...datum.record,
      values: datum.form.fields.map((field) => {
        if (getFieldKind(field.fieldType) === FieldKind.SubForm) {
          return {
            fieldId: field.id,
            value: value[datum.form.id][field.id].map(
              (subValue: RHFieldValue<any>) =>
                Object.entries(subValue).reduce<FieldValue[]>(
                  (acc, [key, v]) => [
                    ...acc,
                    { fieldId: key, value: v as string },
                  ],
                  [],
                ),
            ),
          };
        }
        return {
          fieldId: field.id,
          value: value[datum.form.id][field.id],
        };
      }),
    },
  }));
