import { FieldKind, FormWithRecord, Record } from '../types';
import { getFieldKind } from '../utils';

export const validateRecipientHierarchy = (
  formWithRecordList: FormWithRecord<Record>[],
): void => {
  formWithRecordList.forEach(({ form, record }, idx, arr) => {
    if (form.id !== record.formId)
      throw new Error('Record is not associated with the form');

    if (idx === 0) {
      // assert no key field
      if (form.fields.find((f) => f.key))
        throw new Error(
          `Root recipient form ${form.id} should not have a key field`,
        );
    } else {
      // assert reference key field
      const keyField = form.fields.find(
        (f) => f.key && getFieldKind(f.fieldType) === FieldKind.Reference,
      );

      if (!keyField) throw new Error(`No key field found for form ${form.id}`);

      const { form: prevForm } = arr[idx - 1];

      if (
        keyField.fieldType.reference?.formId !== prevForm.id ||
        keyField.fieldType.reference?.databaseId !== prevForm.databaseId
      )
        throw new Error();
    }
  });
};
