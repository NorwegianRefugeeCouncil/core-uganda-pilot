import { FormDefinition, Record } from 'core-api-client';

export const buildDefaultFormValues = (
  forms: FormDefinition[],
  records: Record[],
): {
  [formId: string]: {
    [fieldId: string]: any;
  };
} =>
  forms.reduce((acc, form) => {
    const record: Record | undefined = records.find(
      (r) => r.formId === form.id,
    );
    if (!record) return acc;
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
