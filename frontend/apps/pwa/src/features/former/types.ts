import { FieldKind, FormType } from 'core-api-client';
import { FieldErrors } from 'react-hook-form';

import { FormField, ValidationForm } from '../../reducers/Former/types';

export type FormerProps = {
  formId: string;
  formType: FormType;
  addField: (kind: FieldKind) => void;
  addOption: (fieldId: string) => void;
  cancelField: (fieldId: string) => void;
  errors: FieldErrors<ValidationForm>;
  fieldOptions?: string[];
  fields: FormField[];
  formName: string;
  invalid: boolean;
  openSubForm: (fieldId: string) => void;
  ownerFormName: string | undefined;
  removeOption: (fieldId: string, index: number) => void;
  saveField: (field: FormField) => void;
  saveForm: () => void;
  selectedFieldId: string | undefined;
  setFieldDescription: (fieldId: string, description: string) => void;
  setFieldIsKey: (fieldId: string, isKey: boolean) => void;
  setFieldName: (fieldId: string, name: string) => void;
  setFieldOption: (fieldId: string, i: number, value: string) => void;
  setFieldReferencedDatabaseId: (fieldId: string, databaseId: string) => void;
  setFieldReferencedFormId: (fieldId: string, formId: string) => void;
  setFieldRequired: (fieldId: string, required: boolean) => void;
  setFormName: (formName: string) => void;
  setSelectedField: (fieldId: string | undefined) => void;
};
