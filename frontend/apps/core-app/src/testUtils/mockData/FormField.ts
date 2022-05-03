import { FieldDefinition, FieldType } from 'core-api-client';

export const makeField = (
  index: number,
  key: boolean,
  required: boolean,
  type: FieldType,
): FieldDefinition => ({
  id: `field-${index}`,
  name: `field-name-${index}`,
  description: `description-${index}`,
  code: '',
  required,
  key,
  fieldType: type,
});
