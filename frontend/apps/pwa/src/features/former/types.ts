import { FieldKind } from 'core-api-client';

import { FormField } from '../../reducers/former';

export type ApiErrorDetails = {
  message: string;
  reason: string;
  field: string;
};

export type FormerProps = {
  formName: string;
  setFormName: (formName: string) => void;
  fields: FormField[];
  fieldOptions?: string[];
  setFieldOption: (fieldId: string, i: number, value: string) => void;
  addOption: (fieldId: string) => void;
  removeOption: (fieldId: string, index: number) => void;
  selectedFieldId: string | undefined;
  setSelectedField: (fieldId: string | undefined) => void;
  addField: (kind: FieldKind) => void;
  setFieldRequired: (fieldId: string, required: boolean) => void;
  setFieldIsKey: (fieldId: string, isKey: boolean) => void;
  setFieldName: (fieldId: string, name: string) => void;
  setFieldDescription: (fieldId: string, description: string) => void;
  setFieldReferencedDatabaseId: (fieldId: string, databaseId: string) => void;
  setFieldReferencedFormId: (fieldId: string, formId: string) => void;
  openSubForm: (fieldId: string) => void;
  saveField: (fieldId: string) => void;
  saveForm: () => void;
  ownerFormName: string | undefined;
  cancelField: (fieldId: string) => void;
  error: ApiErrorDetails[];
};
