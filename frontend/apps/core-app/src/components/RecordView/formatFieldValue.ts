import { FieldDefinition, FieldKind, getFieldKind } from 'core-api-client';

// This method takes a field value and returns a formatted string
// It handles extra cases, like a multi-select field value that is not an array, for safety
export const formatFieldValue = (
  value: string | string[] | null,
  field: FieldDefinition,
): string => {
  const isNull = value === null;
  const isArray = Array.isArray(value);

  switch (getFieldKind(field.fieldType)) {
    case FieldKind.Checkbox:
      return isArray || !value || value === 'false' ? 'No' : 'Yes';
    case FieldKind.Date:
    case FieldKind.Month:
    case FieldKind.Week:
      return isArray || isNull ? '-' : value;
    case FieldKind.Reference:
      // TODO: Lookup reference form?
      return isArray || isNull ? '-' : value;
    case FieldKind.SingleSelect:
      return !isArray
        ? field.fieldType.singleSelect?.options.find((o) => o.id === value)
            ?.name ?? '-'
        : '-';
    case FieldKind.MultiSelect:
      return isArray
        ? value
            .map(
              (v) =>
                field.fieldType.multiSelect?.options.find((o) => o.id === v)
                  ?.name ?? '-',
            )
            .join(', ')
        : field.fieldType.multiSelect?.options.find((o) => o.id === value)
            ?.name ?? '-';
    case FieldKind.SubForm:
      return 'Error: Subform fields should not reach this point';
    default: {
      if (isNull) return '-';
      if (isArray) return value.join(', ');
      return value.toString();
    }
  }
};
