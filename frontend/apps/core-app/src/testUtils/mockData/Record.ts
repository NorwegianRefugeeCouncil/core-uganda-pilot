import { FormDefinition, Record } from 'core-api-client';

export const makeRecord = (index: number, form: FormDefinition): Record => ({
  id: `record-${index}`,
  values: form.fields.map((f) => {
    return {
      value: `value-${f.id}`,
      fieldId: f.id,
    };
  }),
  formId: form.id,
  databaseId: 'database-id',
  ownerId: undefined,
});
