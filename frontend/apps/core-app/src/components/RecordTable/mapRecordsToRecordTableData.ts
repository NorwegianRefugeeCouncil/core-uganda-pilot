import { FieldValue, FormWithRecord, Record } from 'core-api-client';

import { RecordTableEntry } from './types';

export const mapRecordsToRecordTableData = (
  data: FormWithRecord<Record>[][],
): RecordTableEntry[] => {
  return data.reduce((allEntries: RecordTableEntry[], item) => {
    const completeEntry = item.reduce(
      (ce: RecordTableEntry, formWithRecord) => {
        const partialEntry = formWithRecord.record.values.reduce(
          (pe: RecordTableEntry, value: FieldValue) => {
            const field = formWithRecord.form.fields.find((f) => {
              return !f.key && f.id === value.fieldId;
            });
            if (field) return { ...pe, [field?.id]: value.value };
            return pe;
          },
          {},
        );
        return { ...ce, ...partialEntry };
      },
      {},
    );
    return allEntries.concat([completeEntry]);
  }, []);
};
