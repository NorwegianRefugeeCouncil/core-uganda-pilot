import { FieldDefinition, FieldKind, getFieldKind } from 'core-api-client';

export const formatFieldValue = (
  value: string | string[] | null,
  field: FieldDefinition,
): string => {
  if (value === null) return '-';

  const isArray = Array.isArray(value);

  switch (getFieldKind(field.fieldType)) {
    case FieldKind.Checkbox:
      return value === 'true' ? 'Yes' : 'No';
    case FieldKind.Date:
    case FieldKind.Month:
    case FieldKind.Week:
      return !isArray ? value : '-';
    case FieldKind.Reference:
      // TODO: Lookup reference form?
      return !isArray ? value : '-';
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
        : '-';
    case FieldKind.SubForm:
      return 'Error: Subform fields should not reach this point';
    default:
      return isArray ? value.join(', ') : value.toString();
  }
};
