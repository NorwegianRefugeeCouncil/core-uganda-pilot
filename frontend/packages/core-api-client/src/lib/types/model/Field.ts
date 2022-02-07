type EmptyObject = Record<string, never>;

export enum FieldKind {
  Text = 'text',
  MultilineText = 'multilineText',
  Reference = 'reference',
  SubForm = 'subform',
  Date = 'date',
  Quantity = 'quantity',
  SingleSelect = 'singleSelect',
  MultiSelect = 'multiSelect',
  Week = 'week',
  Month = 'month',
  Boolean = 'boolean',
}

export type FieldTypeText = EmptyObject;

export type FieldTypeMultilineText = EmptyObject;

export type FieldTypeDate = EmptyObject;

export type FieldTypeMonth = EmptyObject;

export type FieldTypeQuantity = EmptyObject;

export type FieldTypeWeek = EmptyObject;

export type FieldTypeBoolean = EmptyObject;

export type SelectOption = {
  name: string;
  id: string;
};

export type FieldTypeSingleSelect = {
  options: SelectOption[];
};

export type FieldTypeMultiSelect = {
  options: SelectOption[];
};

export type FieldTypeReference = {
  databaseId: string;
  formId: string;
};

export type FieldTypeSubForm = {
  fields: FieldDefinition[];
};

export type FieldType = {
  text?: FieldTypeText;
  reference?: FieldTypeReference;
  subForm?: FieldTypeSubForm;
  multilineText?: FieldTypeMultilineText;
  date?: FieldTypeDate;
  month?: FieldTypeMonth;
  week?: FieldTypeWeek;
  quantity?: FieldTypeQuantity;
  singleSelect?: FieldTypeSingleSelect;
  multiSelect?: FieldTypeMultiSelect;
  boolean?: FieldTypeBoolean;
};

export type FieldDefinition = {
  id: string;
  code: string;
  name: string;
  description: string;
  required: boolean;
  key: boolean;
  fieldType: FieldType;
};
