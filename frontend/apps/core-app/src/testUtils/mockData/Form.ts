import { FieldDefinition, FormDefinition, FormType } from 'core-api-client';

export const makeForm = (
  index = 0,
  type: FormType = FormType.DefaultFormType,
  fields: FieldDefinition[] = [],
): FormDefinition => ({
  id: `form${index}`,
  code: '',
  name: `form ${index}`,
  databaseId: 'databaseId',
  folderId: 'folderId',
  formType: type,
  fields,
});
