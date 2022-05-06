import { FieldValue, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { RecipientListTableEntry } from './types';

export const mapRecordsToRecipientTableData = (
  data: FormWithRecord<Recipient>[][],
): RecipientListTableEntry[] => {
  return data.reduce((allEntries: RecipientListTableEntry[], item) => {
    const completeEntry = item.reduce(
      (ce: RecipientListTableEntry, formWithRecord) => {
        const partialEntry = formWithRecord.record.values.reduce(
          (pe: RecipientListTableEntry, value: FieldValue) => {
            const field = formWithRecord.form.fields.find((f) => {
              return f.id === value.fieldId;
            });
            if (field) return { ...pe, [field?.id]: value.value };
            return pe;
          },
          {},
        );
        return { ...ce, ...partialEntry };
      },
      { recordId: item[item.length - 1].record.id },
    );
    return allEntries.concat([completeEntry]);
  }, []);
};
