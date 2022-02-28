import { FieldDefinition } from './Field';

export enum FormType {
  DefaultFormType = 'default',
  RecipientFormType = 'recipient',
}

export type FormDefinition = {
  id: string;
  code: string;
  databaseId: string;
  folderId: string;
  name: string;
  formType: FormType;
  fields: FieldDefinition[];
};

export type FormDefinitionList = {
  items: FormDefinition[];
};
