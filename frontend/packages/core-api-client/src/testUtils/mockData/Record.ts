import { FormDefinition, Record } from '../../types';

export const makeRecord = (index: number, form: FormDefinition): Record => ({
  id: `record${index}`,
  values: form.fields.map((f) => {
    return {
      value: `value-${f.id}`,
      fieldId: f.id,
    };
  }),
  formId: form.id,
  databaseId: 'databaseId',
  ownerId: undefined,
});
