import { FieldKind, FormType } from 'core-api-client';

import { customValidation, validationConstants } from '../validation';
import { Form } from '../../../reducers/Former/types';

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

      describe('number of fields', () => {
        describe('success', () => {
          it('should return an empty array when the amount of fields is valid', () => {
            const result = customValidation.form({
              ...baseForm,
              fields: Array(1),
            });
            expect(result).toEqual([]);
          });
        });

        describe('failure', () => {
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
        });
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

      describe('valid', () => {
        it('should return an empty array for a valid field', () => {
          const result = customValidation.selectedField(baseField);
          expect(result).toEqual([]);
        });
      });

      describe('singleSelect, amount of options', () => {
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
      });

      describe('subforms, required', () => {
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

    describe('recipient forms', () => {
      let baseForm: Form;
      beforeEach(() => {
        baseForm = {
          name: 'testform',
          formId: 'formId',
          fields: [],
          formType: FormType.RecipientFormType,
          isRootForm: true,
        };
      });

      describe('success', () => {
        it('should not return an error for a valid recipient form', () => {
          baseForm.fields = [
            {
              id: 'id',
              fieldType: FieldKind.Reference,
              options: [],
              required: true,
              key: true,
              name: 'name',
              description: 'description',
              code: 'code',
              subFormId: undefined,
              referencedDatabaseId: undefined,
              referencedFormId: undefined,
            },
          ];
          const result = customValidation.form(baseForm);
          expect(result).toEqual([]);
        });
      });
      describe('failure', () => {
        it('should return an error if form has no key field', () => {
          baseForm.fields = [
            {
              id: 'id',
              fieldType: FieldKind.Reference,
              options: [],
              required: true,
              key: false,
              name: 'name',
              description: 'description',
              code: 'code',
              subFormId: undefined,
              referencedDatabaseId: undefined,
              referencedFormId: undefined,
            },
          ];
          const result = customValidation.form(baseForm);
          expect(result).toEqual([
            {
              field: 'fields',
              message: 'Form needs to have exactly 1 key field',
            },
          ]);
        });
        it('should return an error if form has more than one key field', () => {
          baseForm.fields = [
            {
              id: 'id1',
              fieldType: FieldKind.Reference,
              options: [],
              required: true,
              key: true,
              name: 'name',
              description: 'description',
              code: 'code',
              subFormId: undefined,
              referencedDatabaseId: undefined,
              referencedFormId: undefined,
            },
            {
              id: 'id2',
              fieldType: FieldKind.Reference,
              options: [],
              required: true,
              key: true,
              name: 'name',
              description: 'description',
              code: 'code',
              subFormId: undefined,
              referencedDatabaseId: undefined,
              referencedFormId: undefined,
            },
          ];
          const result = customValidation.form(baseForm);
          expect(result).toEqual([
            {
              field: 'fields',
              message: 'Form needs to have exactly 1 key field',
            },
          ]);
        });
        it('should return an error if the key field is not a reference', () => {
          baseForm.fields = [
            {
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
            },
          ];
          const result = customValidation.form(baseForm);
          expect(result).toEqual([
            {
              field: 'fields',
              message: 'Form needs to have exactly 1 key field',
            },
          ]);
        });
      });
    });
  });
});
