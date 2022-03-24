import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FieldValue } from 'react-hook-form';

export const toReactHookForm = (
  data: FormWithRecord<Recipient>[],
): FieldValue<any> =>
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
  value: FieldValue<any>,
): FormWithRecord<Recipient>[] =>
  originalData.map((datum) => ({
    form: datum.form,
    record: {
      ...datum.record,
      values: datum.form.fields.map((field) => ({
        fieldId: field.id,
        value: value[datum.form.id][field.id],
      })),
    },
  }));
