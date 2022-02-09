import { FieldKind, FormType } from 'core-api-client';
import { FieldPath } from 'react-hook-form';

import { customValidation, validationConstants } from '../validation';
import { ValidationForm } from '../../../reducers/Former/types';

describe('validation', () => {
  describe('customValidation', () => {
    describe('form', () => {
      const baseForm = {
        name: 'testform',
        formId: 'formId',
        fields: [],
        formType: FormType.DefaultFormType,
        isRootForm: true,
      };

      it('should create an error if the form does not have enough fields', () => {
        const result = customValidation.form({ ...baseForm });
        expect(result).toEqual([
          {
            field: 'fields',
            message: `Form needs to have at least ${validationConstants.fields.min} field`,
          },
        ]);
      });

      it('should create an error if the form has too many fields', () => {
        const result = customValidation.form({
          ...baseForm,
          fields: Array(validationConstants.fields.max + 1),
        });
        expect(result).toEqual([
          {
            field: 'fields',
            message: `Form can have at most ${validationConstants.fields.max} fields`,
          },
        ]);
      });

      it('should return an empty array when the amount of fields is valid', () => {
        const result = customValidation.form({ ...baseForm, fields: Array(1) });
        expect(result).toEqual([]);
      });
    });

    describe('selectedField', () => {
      const baseField = {
        id: 'id',
        fieldType: FieldKind.Text,
        options: [],
        required: true,
        key: false,
        name: 'name',
        description: 'description',
        code: 'code',
        subFormId: undefined,
        referencedDatabaseId: undefined,
        referencedFormId: undefined,
      };

      it('should return an empty array for a valid field', () => {
        const result = customValidation.selectedField(baseField);
        expect(result).toEqual([]);
      });

      it('should return an error if a single select does not have enough options', () => {
        const result = customValidation.selectedField({
          ...baseField,
          fieldType: FieldKind.SingleSelect,
          options: [],
        });
        expect(result).toEqual([
          {
            field: 'selectedField.fieldType.singleSelect.options',
            message: `At least ${validationConstants.options.min} options are required`,
          },
        ]);
      });
      it('should return an error if a single select has too many options', () => {
        const result = customValidation.selectedField({
          ...baseField,
          fieldType: FieldKind.SingleSelect,
          options: Array(validationConstants.options.max + 1),
        });
        expect(result).toEqual([
          {
            field: 'selectedField.fieldType.singleSelect.options',
            message: `At most ${validationConstants.options.max} options are allowed`,
          },
        ]);
      });
      it('should return an error if a subform is required', () => {
        const result = customValidation.selectedField({
          ...baseField,
          fieldType: FieldKind.SubForm,
          required: true,
        });
        expect(result).toEqual([
          {
            field: 'selectedField.fieldType.subForm',
            message: 'Subforms cannot be required',
          },
        ]);
      });
    });
  });
});
