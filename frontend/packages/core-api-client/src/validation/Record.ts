import { RegisterOptions } from 'react-hook-form/dist/types/validator';

import { FieldDefinition, FieldKind } from '../types';

const getBaseFieldValidationRules = (
  field: FieldDefinition,
): RegisterOptions => ({
  required: {
    value: field.required,
    message: `${field.name} is a required field`,
  },
});

export const formValidationRules: {
  field: Record<FieldKind, (field: FieldDefinition) => RegisterOptions>;
} = {
  field: {
    [FieldKind.Text]: (field) => ({ ...getBaseFieldValidationRules(field) }),
    [FieldKind.MultilineText]: (field) => ({
      ...getBaseFieldValidationRules(field),
    }),
    [FieldKind.Quantity]: (field) => ({
      ...getBaseFieldValidationRules(field),
      valueAsNumber: true,
    }),
    [FieldKind.Checkbox]: (field) => ({
      ...getBaseFieldValidationRules(field),
    }),
    [FieldKind.Date]: (field) => ({
      ...getBaseFieldValidationRules(field),
      valueAsDate: true,
    }),
    [FieldKind.Month]: (field) => ({
      ...getBaseFieldValidationRules(field),
      valueAsDate: true,
      pattern: {
        value: /^(?:19|20|21)\d{2}-[01]\d$/,
        message: 'Invalid month format',
      },
    }),
    [FieldKind.Week]: (field) => ({
      ...getBaseFieldValidationRules(field),
      valueAsDate: true,
      pattern: {
        value: /^(?:19|20|21)\d{2}-W[0-5]\d$/,
        message: 'Invalid week format',
      },
    }),
    [FieldKind.SingleSelect]: (field) => ({
      ...getBaseFieldValidationRules(field),
    }),
    [FieldKind.MultiSelect]: (field) => ({
      ...getBaseFieldValidationRules(field),
    }),
    [FieldKind.Reference]: (field) => ({
      ...getBaseFieldValidationRules(field),
    }),
    [FieldKind.SubForm]: (field) => ({ ...getBaseFieldValidationRules(field) }),
  },
};
