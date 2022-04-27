import { FieldValue, FormWithRecord, Record } from 'core-api-client';

import { RecipientListTableEntry } from './types';

export const mapRecordsToRecordTableData = (
  data: FormWithRecord<Record>[][],
): RecipientListTableEntry[] => {
  return data.reduce((allEntries: RecipientListTableEntry[], item) => {
    const completeEntry = item.reduce(
      (ce: RecipientListTableEntry, formWithRecord) => {
        const partialEntry = formWithRecord.record.values.reduce(
          (pe: RecipientListTableEntry, value: FieldValue) => {
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
