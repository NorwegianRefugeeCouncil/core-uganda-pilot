import { FieldValue, FormDefinition, Record } from 'core-api-client';

export const mapRecordsToTableData = (
  records: Record[],
  form: FormDefinition,
): any[] =>
  records.map((record) => {
    return record.values.reduce((acc, value: FieldValue) => {
      const field = form.fields.find((f) => {
        return f.id === value.fieldId;
      });
      if (field) return { ...acc, [field?.id]: value.value };
      return acc;
    }, {});
  });
