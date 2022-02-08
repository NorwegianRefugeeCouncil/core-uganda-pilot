import { FieldKind, FormType, SelectOption } from 'core-api-client';
import { EntityState } from '@reduxjs/toolkit';
import {
  FieldTypeCheckbox,
  FieldTypeDate,
  FieldTypeMonth,
  FieldTypeMultilineText,
  FieldTypeMultiSelect,
  FieldTypeQuantity,
  FieldTypeReference,
  FieldTypeSingleSelect,
  FieldTypeText,
  FieldTypeWeek,
} from 'core-api-client/src';

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

export type FieldDefinitionNC = {
  id: string;
  code: string;
  name: string;
  description: string;
  required: boolean;
  key: boolean;
  fieldType: FieldTypeNC;
};

export interface FieldTypeNC {
  checkbox?: FieldTypeCheckbox;
  date?: FieldTypeDate;
  month?: FieldTypeMonth;
  multiSelect?: FieldTypeMultiSelect;
  multilineText?: FieldTypeMultilineText;
  quantity?: FieldTypeQuantity;
  reference?: FieldTypeReference;
  singleSelect?: FieldTypeSingleSelect;
  subForm?: FieldTypeSubFormNC;
  text?: FieldTypeText;
  week?: FieldTypeWeek;
}

export type FieldTypeSubFormNC = {
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

export type ValidationForm = Form & { selectedField?: FieldDefinitionNC };
