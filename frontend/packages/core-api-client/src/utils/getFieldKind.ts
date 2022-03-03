import { FieldKind, FieldType } from '../types';

export function getFieldKind(fieldType: FieldType): FieldKind {
  if (fieldType.text) {
    return FieldKind.Text;
  }
  if (fieldType.multilineText) {
    return FieldKind.MultilineText;
  }
  if (fieldType.reference) {
    return FieldKind.Reference;
  }
  if (fieldType.subForm) {
    return FieldKind.SubForm;
  }
  if (fieldType.date) {
    return FieldKind.Date;
  }
  if (fieldType.quantity) {
    return FieldKind.Quantity;
  }
  if (fieldType.singleSelect) {
    return FieldKind.SingleSelect;
  }
  if (fieldType.multiSelect) {
    return FieldKind.MultiSelect;
  }
  if (fieldType.week) {
    return FieldKind.Week;
  }
  if (fieldType.month) {
    return FieldKind.Month;
  }
  if (fieldType.checkbox) {
    return FieldKind.Checkbox;
  }
  throw new Error('unknown field kind');
}
