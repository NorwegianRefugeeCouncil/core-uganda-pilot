import {
  FieldKind,
  FormDefinition,
  getFieldKind,
  Record as CoreRecord,
} from 'core-api-client';

import { formatFieldValue } from './formatFieldValue';
import {
  NormalisedBasicField,
  NormalisedFieldValue,
  NormalisedSubFormField,
} from './RecordView.types';

export const normaliseFieldValues = (
  form: FormDefinition,
  record: CoreRecord,
): NormalisedFieldValue[] =>
  record.values.reduce<NormalisedFieldValue[]>((acc, cur) => {
    const field = form.fields.find((f) => f.id === cur.fieldId);
    if (!field) return acc;
    const fieldType = getFieldKind(field.fieldType);

    if (fieldType === FieldKind.SubForm) {
      if (!Array.isArray(cur.value)) return acc;
      const item: NormalisedSubFormField = {
        key: field.key,
        header: field.name,
        fieldType: FieldKind.SubForm,
        columns:
          field.fieldType.subForm?.fields.map((f) => ({
            Header: f.name,
            accessor: f.id,
          })) ?? [],
        data: cur.value.map((subRecord) => {
          if (!Array.isArray(subRecord)) return {};
          return subRecord.reduce<Record<string, string>>((a, c) => {
            const subField = field.fieldType.subForm?.fields.find(
              (f) => f.id === c.fieldId,
            );
            if (!subField) throw new Error('subField not found');
            return {
              ...a,
              [c.fieldId]: formatFieldValue(
                c.value as string | string[] | null,
                subField,
              ),
            };
          }, {});
        }),
      };

      return [...acc, item];
    }

    const item: NormalisedBasicField = {
      key: field.key,
      label: field.name,
      value: cur.value as string | string[] | null,
      fieldType,
      formattedValue: formatFieldValue(
        cur.value as string | string[] | null,
        field,
      ),
    };

    return [...acc, item];
  }, []);
