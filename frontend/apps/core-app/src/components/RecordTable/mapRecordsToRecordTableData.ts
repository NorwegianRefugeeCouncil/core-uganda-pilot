import { FieldValue, FormWithRecords, Record } from 'core-api-client';

export const mapRecordsToRecordTableData = ({
  records,
  form,
}: FormWithRecords<Record>): any[] =>
  records.map((record) => {
    return record.values.reduce((acc, value: FieldValue) => {
      const field = form.fields.find((f) => {
        return f.id === value.fieldId;
      });
      if (field) return { ...acc, [field?.id]: value.value };
      return acc;
    }, {});
  });
