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
      if (!Array.isArray(cur.value)) return acc;
      const values: NormalisedSubFormFieldValue[][] = cur.value.map(
        (subRecord) => {
          if (!Array.isArray(subRecord)) return [];
          return subRecord.map((subFieldValue) => {
            const subField = field.fieldType.subForm?.fields.find(
              (f) => f.id === subFieldValue.fieldId,
            );
            if (!subField) throw new Error('subField not found');
            const subFieldType = getFieldKind(subField.fieldType);
            if (subFieldType === FieldKind.SubForm)
              throw new Error('subField is a subform');
            return {
              value: subFieldValue.value as string | string[] | null,
              formattedValue: formatFieldValue(
                subFieldValue.value as string | string[] | null,
                subField,
              ),
              fieldType: subFieldType,
            };
          });
        },
      );

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
