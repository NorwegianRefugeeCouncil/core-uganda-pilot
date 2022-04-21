import { FieldKind, FormWithRecord, Record } from '../types';
import { getFieldKind } from '../utils';

export const validateRecipientHierarch = (
  formWithRecordList: FormWithRecord<Record>[],
): void => {
  formWithRecordList.forEach(({ form, record }, idx, arr) => {
    if (form.id !== record.formId) throw new Error();

    if (idx === 0) {
      // assert no key field
      if (form.fields.find((f) => f.key)) throw new Error();
    } else {
      // assert reference key field
      const keyField = form.fields.find(
        (f) => f.key && getFieldKind(f.fieldType) === FieldKind.Reference,
      );

      if (!keyField) throw new Error();

      const { form: prevForm } = arr[idx - 1];

      if (
        keyField.fieldType.reference?.formId !== prevForm.id ||
        keyField.fieldType.reference?.databaseId !== prevForm.databaseId
      )
        throw new Error();
    }
  });
};
