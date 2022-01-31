import { FieldKind, FormType, SelectOption } from 'core-api-client';
import { EntityState } from '@reduxjs/toolkit';

import { ErrorMessage } from '../../types/errors';

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
  errors: ErrorMessage | undefined;
}

export interface Form {
  // name of the form
  name: string;
  // the unique id of the form
  formId: string;
  type: FormType;
  // records the record values
  fields: FormField[];
  isRootForm: boolean;
  errors: ErrorMessage | undefined;
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
