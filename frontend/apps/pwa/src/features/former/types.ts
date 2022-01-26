import { FieldKind, FormType } from 'core-api-client';
import { FieldErrors } from 'react-hook-form';

import { Form, FormField } from '../../reducers/Former/types';

export type FormerProps = {
  formId: string;
  formType: FormType;
  addField: (kind: FieldKind) => void;
  addOption: (fieldId: string) => void;
  cancelField: (fieldId: string) => void;
  errors: FieldErrors<Form>;
  fieldOptions?: string[];
  fields: FormField[];
  formName: string;
  openSubForm: (fieldId: string) => void;
  ownerFormName: string | undefined;
  register: any;
  removeOption: (fieldId: string, index: number) => void;
  saveField: (fieldId: string) => void;
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
  invalid: boolean;
};
