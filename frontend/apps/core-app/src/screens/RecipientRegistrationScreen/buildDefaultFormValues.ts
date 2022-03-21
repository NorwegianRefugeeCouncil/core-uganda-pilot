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
    if (!record)
      throw new Error(
        `Cannot find record for form ${form.id} when building default react-hook-form values`,
      );
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
