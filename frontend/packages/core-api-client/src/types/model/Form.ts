import { FieldDefinition } from './Field';
import { Record } from './Record';

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

export type FormWithRecord<T extends Record> = {
  form: FormDefinition;
  record: T;
};
