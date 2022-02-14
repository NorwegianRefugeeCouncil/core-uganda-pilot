import {
  FieldKind,
  FieldTypeCheckbox,
  FieldTypeDate,
  FieldTypeMonth,
  FieldTypeMultiSelect,
  FieldTypeMultilineText,
  FieldTypeQuantity,
  FieldTypeReference,
  FieldTypeSingleSelect,
  FieldTypeText,
  FieldTypeWeek,
  FormType,
  SelectOption,
} from 'core-api-client';
import { EntityState } from '@reduxjs/toolkit';

export interface FormField {
  id: string;
  fieldType: FieldKind;
  options: SelectOption[];
  required: boolean;
  key: boolean;
  name: string;
  description: string;
  code: string;
  subFormId: string | undefined;
  referencedDatabaseId: string | undefined;
  referencedFormId: string | undefined;
}

export type FieldDefinitionNonCircular = {
  id: string;
  code: string;
  name: string;
  description: string;
  required: boolean;
  key: boolean;
  fieldType: FieldTypeNonCircular;
};

export interface FieldTypeNonCircular {
  checkbox?: FieldTypeCheckbox;
  date?: FieldTypeDate;
  month?: FieldTypeMonth;
  multiSelect?: FieldTypeMultiSelect;
  multilineText?: FieldTypeMultilineText;
  quantity?: FieldTypeQuantity;
  reference?: FieldTypeReference;
  singleSelect?: FieldTypeSingleSelect;
  subForm?: FieldTypeSubFormNonCircular;
  text?: FieldTypeText;
  week?: FieldTypeWeek;
}

export type FieldTypeSubFormNonCircular = {
  id: string;
};

export interface Form {
  // name of the form
  name: string;
  // the unique id of the form
  formId: string;
  formType: FormType;
  // records the record values
  fields: FormField[];
  isRootForm: boolean;
}

export interface FormerState extends EntityState<Form> {
  selectedFormId: string;
  selectedFieldId: string | undefined;
  selectedDatabaseId: string | undefined;
  selectedFolderId: string | undefined;
  savePending: boolean;
  saveSuccess: boolean;
  saveError: any;
}

export type ValidationForm = Form & {
  selectedField?: FieldDefinitionNonCircular;
};
