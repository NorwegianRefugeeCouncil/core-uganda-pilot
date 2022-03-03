import {
  FieldKind,
  FormDefinition,
  getFieldKind,
  Record,
} from 'core-api-client';

import { formatFieldValue } from './formatFieldValue';
import {
  NormalisedBasicField,
  NormalisedFieldValue,
  NormalisedSubFormField,
  NormalisedSubFormFieldValue,
} from './RecordView.types';

export const normaliseFieldValues = (
  form: FormDefinition,
  record: Record,
): NormalisedFieldValue[] =>
  record.values.reduce<NormalisedFieldValue[]>((acc, cur) => {
    const field = form.fields.find((f) => f.id === cur.fieldId);
    if (!field) return acc;
    const fieldType = getFieldKind(field.fieldType);

    if (fieldType === FieldKind.SubForm) {
      const values: NormalisedSubFormFieldValue[][] =
        field.fieldType.subForm?.fields.map((subField, i) => {
          if (!Array.isArray(cur.value)) return [];
          return cur.value.map((v) => {
            if (!Array.isArray(v)) throw new Error();
            const vv = v[i].value as string | string[] | null;
            const subFieldType = getFieldKind(subField.fieldType);
            if (subFieldType === FieldKind.SubForm) throw new Error();
            return {
              value: vv,
              fieldType: subFieldType,
              formattedValue: formatFieldValue(vv, subField),
            };
          });
        }) ?? [];

      const labels =
        field.fieldType.subForm?.fields.map((subField) => subField.name) ?? [];

      const item: NormalisedSubFormField = {
        fieldType,
        header: field.name,
        labels,
        values,
      };

      return [...acc, item];
    }

    const item: NormalisedBasicField = {
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
