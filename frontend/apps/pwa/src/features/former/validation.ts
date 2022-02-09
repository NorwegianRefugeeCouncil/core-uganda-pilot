import { FieldPath } from 'react-hook-form';
import { FieldKind } from 'core-api-client';

import { Form, FormField, ValidationForm } from '../../reducers/Former/types';

const validationConstants = {
  name: {
    minLength: 3,
    maxLength: 128,
    pattern: /[A-Za-z0-9]/,
  },
  options: {
    min: 2,
    max: 60,
    name: {
      minLength: 2,
      maxLength: 128,
    },
  },
  fields: {
    min: 1,
    max: 100,
  },
};

export const registeredValidation = {
  name: {
    minLength: {
      value: validationConstants.name.minLength,
      message: `Form name needs to be at least ${validationConstants.name.minLength} characters long`,
    },
    maxLength: {
      value: validationConstants.name.maxLength,
      message: `Form name can be at most ${validationConstants.name.maxLength} characters long`,
    },
    pattern: {
      value: validationConstants.name.pattern,
      message: 'Form name contains invalid characters',
    },
    required: { value: true, message: 'Form name is a required field' },
  },
  selectedField: {
    name: {
      minLength: {
        value: validationConstants.name.minLength,
        message: `Field name needs to be at least ${validationConstants.name.minLength} characters long`,
      },
      maxLength: {
        value: validationConstants.name.maxLength,
        message: `Field name can be at most ${validationConstants.name.maxLength} characters long`,
      },
      required: { value: true, message: 'Field name is required' },
    },
    fieldType: {
      reference: {
        databaseId: {
          required: { value: true, message: 'Data base is required' },
        },
        formId: {
          required: { value: true, message: 'Form is required' },
        },
      },
      select: {
        option: {
          name: {
            minLength: {
              value: validationConstants.name.minLength,
              message: `Option name needs to be at least ${validationConstants.name.minLength} characters long`,
            },
            maxLength: {
              value: validationConstants.name.maxLength,
              message: `Option name can be at most ${validationConstants.name.maxLength} characters long`,
            },
            required: { value: true, message: 'Option name is required' },
          },
        },
      },
    },
  },
};

type CustomError = {
  field: FieldPath<ValidationForm>;
  message: string;
};

export const customValidation = {
  form: (form: Form): CustomError[] => {
    const errors = [];
    if (form.fields.length < validationConstants.fields.min) {
      errors.push({
        field: 'fields' as FieldPath<ValidationForm>,
        message: `Form needs to have at least ${validationConstants.fields.min} field`,
      });
    }
    if (form.fields.length > validationConstants.fields.max) {
      errors.push({
        field: 'fields' as FieldPath<ValidationForm>,
        message: `Form can have at most ${validationConstants.fields.max} fields`,
      });
    }
    return errors;
  },
  selectedField: (field: FormField): CustomError[] => {
    const errors = [];
    if (
      field.fieldType === FieldKind.SingleSelect ||
      field.fieldType === FieldKind.MultiSelect
    ) {
      if (field.options.length < validationConstants.options.min) {
        errors.push({
          field:
            `selectedField.fieldType.${field.fieldType}.options` as FieldPath<ValidationForm>,
          message: `At least ${validationConstants.options.min} options are required`,
        });
      }
      if (field.options.length > validationConstants.options.max) {
        errors.push({
          field:
            `selectedField.fieldType.${field.fieldType}.options` as FieldPath<ValidationForm>,
          message: `At most ${validationConstants.options.max} options are allowed`,
        });
      }
    }
    if (field.fieldType === FieldKind.SubForm) {
      if (field.required) {
        errors.push({
          field: 'selectedField.fieldType.subForm' as FieldPath<ValidationForm>,
          message: 'Subforms cannot be required',
        });
      }
    }
    return errors;
  },
};
