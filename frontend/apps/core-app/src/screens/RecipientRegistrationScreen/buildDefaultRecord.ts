import {
  FieldKind,
  FormDefinition,
  getFieldKind,
  Record,
} from 'core-api-client';

export const buildDefaultRecord = (form: FormDefinition): Record => ({
  id: '', // We are creating records so there shouldn't be an id
  databaseId: form.databaseId,
  formId: form.id,
  ownerId: undefined,
  values: form.fields.map((field) => {
    const fieldType = getFieldKind(field.fieldType);
    switch (fieldType) {
      case FieldKind.Text:
      case FieldKind.MultilineText:
        return {
          fieldId: field.id,
          value: '',
        };
      case FieldKind.Reference:
      case FieldKind.Date:
      case FieldKind.Month:
      case FieldKind.Week:
      case FieldKind.SingleSelect:
        return {
          fieldId: field.id,
          value: null,
        };
      case FieldKind.MultiSelect:
        return {
          fieldId: field.id,
          value: [],
        };
      case FieldKind.Checkbox:
        return {
          fieldId: field.id,
          value: 'false',
        };
      case FieldKind.SubForm:
        return {
          fieldId: field.id,
          value: [],
        };
      default:
        return {
          fieldId: field.id,
          value: '',
        };
    }
  }),
});
