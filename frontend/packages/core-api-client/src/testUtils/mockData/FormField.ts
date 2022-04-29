import { FieldDefinition, FieldType } from '../../types';

export const makeField = (
  index: number,
  key: boolean,
  required: boolean,
  type: FieldType,
): FieldDefinition => ({
  id: `field${index}`,
  name: `field ${index}`,
  description: `description ${index}`,
  code: '',
  required,
  key,
  fieldType: type,
});
