import {
  Record,
  FieldKind,
  FormWithRecord,
  getFieldKind,
  FormType,
  FieldValue,
} from 'core-api-client';
import { FieldValue as RHFFieldValue } from 'react-hook-form';

export const toReactHookForm = (
  data: FormWithRecord<Record>[],
): RHFFieldValue<any> =>
  data.reduce((acc, { form, record }) => {
    return {
      ...acc,
      [form.id]: record.values.reduce((innerAcc, fieldValue) => {
        const field = form.fields.find((f) => f.id === fieldValue.fieldId);
        if (!field) return innerAcc;
        if (getFieldKind(field.fieldType) === FieldKind.SubForm) {
          return {
            ...innerAcc,
            [fieldValue.fieldId]: (fieldValue.value as FieldValue[][]).map(
              (v) =>
                v.reduce(
                  (a, { fieldId, value }) => ({
                    ...a,
                    [fieldId]: value,
                  }),
                  {},
                ),
            ),
          };
        }
        return {
          ...innerAcc,
          [fieldValue.fieldId]: fieldValue.value,
        };
      }, {}),
    };
  }, {});

export const fromReactHookForm = (
  originalData: FormWithRecord<Record>[],
  value: RHFFieldValue<any>,
): FormWithRecord<Record>[] =>
  originalData.map((datum) => ({
    form: datum.form,
    record: {
      ...datum.record,
      values: datum.form.fields.map((field) => {
        if (getFieldKind(field.fieldType) === FieldKind.SubForm) {
          const v = value[datum.form.id][field.id].map(
            (vv) =>
              fromReactHookForm(
                [
                  {
                    form: {
                      id: field.id,
                      databaseId: datum.form.databaseId,
                      folderId: datum.form.folderId,
                      name: '',
                      code: '',
                      formType: FormType.DefaultFormType,
                      fields: field.fieldType.subForm?.fields ?? [],
                    },
                    record: {
                      id: '',
                      databaseId: datum.form.databaseId,
                      formId: field.id,
                      ownerId: datum.record.id,
                      values: [],
                    },
                  },
                ],
                { [field.id]: vv },
              )[0].record.values,
          );
          return {
            fieldId: field.id,
            value: v,
          };
        }
        return {
          fieldId: field.id,
          value: value[datum.form.id][field.id],
        };
      }),
    },
  }));
