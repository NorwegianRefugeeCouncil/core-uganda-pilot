import * as React from 'react';
import { FormDefinition, FormType, Record } from 'core-api-client';

import { buildDefaultRecord } from './buildDefaultRecord';
import { RecipientRegistrationScreenComponent } from './RecipientRegistrationScreen.component';

const DATABASE_ID = '9eab5b48-41f1-411a-8cf0-390dfee8ca37';
const FORM_ID = 'e4e56318-2147-49ad-901e-a8c3f451079b';

const makeForm = (id: string, name: string): FormDefinition => ({
  id,
  code: '',
  databaseId: DATABASE_ID,
  folderId: '',
  formType: FormType.DefaultFormType,
  name,
  fields: [
    {
      id: 'text-field',
      name: 'Text Field',
      description: "I'm a text field",
      code: '',
      required: false,
      key: false,
      fieldType: { text: {} },
    },
    {
      id: 'multi-line-text-field',
      name: 'Multi-line Text Field',
      description: "I'm a multi-line text field",
      code: '',
      required: false,
      key: false,
      fieldType: { multilineText: {} },
    },
    {
      id: 'quantity-field',
      name: 'Quantity Field',
      description: "I'm a quantity field",
      code: '',
      required: false,
      key: false,
      fieldType: { quantity: {} },
    },
    {
      id: 'checkbox-field',
      name: 'Checkbox Field',
      description: "I'm a checkbox field",
      code: '',
      required: false,
      key: false,
      fieldType: { checkbox: {} },
    },
    {
      id: 'date-field',
      name: 'Date Field',
      description: "I'm a date field",
      code: '',
      required: false,
      key: false,
      fieldType: { date: {} },
    },
    {
      id: 'month-field',
      name: 'Month Field',
      description: "I'm a month field",
      code: '',
      required: false,
      key: false,
      fieldType: { month: {} },
    },
    {
      id: 'week-field',
      name: 'Week Field',
      description: "I'm a week field",
      code: '',
      required: false,
      key: false,
      fieldType: { week: {} },
    },
    {
      id: 'single-select-field',
      name: 'Single Select Field',
      description: "I'm a single select field",
      code: '',
      required: false,
      key: false,
      fieldType: {
        singleSelect: {
          options: [
            { id: 'option-1', name: 'Option 1' },
            { id: 'option-2', name: 'Option 2' },
            { id: 'option-3', name: 'Option 3' },
          ],
        },
      },
    },
    {
      id: 'multi-select-field',
      name: 'Multi Select Field',
      description: "I'm a multi select field",
      code: '',
      required: false,
      key: false,
      fieldType: {
        multiSelect: {
          options: [
            { id: 'option-1', name: 'Option 1' },
            { id: 'option-2', name: 'Option 2' },
            { id: 'option-3', name: 'Option 3' },
          ],
        },
      },
    },
    {
      id: 'reference-field',
      name: 'Reference Field',
      description: "I'm a reference field",
      code: '',
      required: false,
      key: false,
      fieldType: {
        reference: {
          databaseId: DATABASE_ID,
          formId: FORM_ID,
        },
      },
    },
    {
      id: 'sub-form-field',
      name: 'Sub Form Field',
      description: "I'm a sub form field",
      code: '',
      required: false,
      key: false,
      fieldType: {
        subForm: {
          fields: [
            {
              id: 'sub-form-text-field',
              name: 'Sub Form Text Field',
              description: "I'm a sub form text field",
              code: '',
              required: false,
              key: false,
              fieldType: { text: {} },
            },
            {
              id: 'sub-form-multi-line-text-field',
              name: 'Sub Form Multi-line Text Field',
              description: "I'm a sub form multi-line text field",
              code: '',
              required: false,
              key: false,
              fieldType: { multilineText: {} },
            },
            {
              id: 'sub-form-quantity-field',
              name: 'Sub Form Quantity Field',
              description: "I'm a sub form quantity field",
              code: '',
              required: false,
              key: false,
              fieldType: { quantity: {} },
            },
            {
              id: 'sub-form-checkbox-field',
              name: 'Sub Form Checkbox Field',
              description: "I'm a sub form checkbox field",
              code: '',
              required: false,
              key: false,
              fieldType: { checkbox: {} },
            },
            {
              id: 'sub-form-date-field',
              name: 'Sub Form Date Field',
              description: "I'm a sub form date field",
              code: '',
              required: false,
              key: false,
              fieldType: { date: {} },
            },
            {
              id: 'sub-form-month-field',
              name: 'Sub Form Month Field',
              description: "I'm a sub form month field",
              code: '',
              required: false,
              key: false,
              fieldType: { month: {} },
            },
            {
              id: 'sub-form-week-field',
              name: 'Sub Form Week Field',
              description: "I'm a sub form week field",
              code: '',
              required: false,
              key: false,
              fieldType: { week: {} },
            },
            {
              id: 'sub-form-single-select-field',
              name: 'Sub Form Single Select Field',
              description: "I'm a sub form single select field",
              code: '',
              required: false,
              key: false,
              fieldType: {
                singleSelect: {
                  options: [
                    { id: 'option-1', name: 'Option 1' },
                    { id: 'option-2', name: 'Option 2' },
                    { id: 'option-3', name: 'Option 3' },
                  ],
                },
              },
            },
            {
              id: 'sub-form-multi-select-field',
              name: 'Sub Form Multi Select Field',
              description: "I'm a sub form multi select field",
              code: '',
              required: false,
              key: false,
              fieldType: {
                multiSelect: {
                  options: [
                    { id: 'option-1', name: 'Option 1' },
                    { id: 'option-2', name: 'Option 2' },
                    { id: 'option-3', name: 'Option 3' },
                  ],
                },
              },
            },
            {
              id: 'sub-form-reference-field',
              name: 'Sub Form Reference Field',
              description: "I'm a sub form reference field",
              code: '',
              required: false,
              key: false,
              fieldType: {
                reference: {
                  databaseId: DATABASE_ID,
                  formId: FORM_ID,
                },
              },
            },
          ],
        },
      },
    },
  ],
});

const getForms = async (): Promise<FormDefinition[]> =>
  Promise.resolve([makeForm('form-1', 'Form 1'), makeForm('form-2', 'Form 2')]);

export const RecipientRegistrationScreenContainer: React.FC = () => {
  const [mode, setMode] = React.useState<'register' | 'review'>('register');
  const [forms, setForms] = React.useState<FormDefinition[]>([]);
  const [records, setRecords] = React.useState<Record[]>([]);

  React.useEffect(() => {
    (async () => {
      const formsResponse = await getForms();
      setForms(formsResponse);
    })();
  }, []);

  React.useEffect(() => {
    setRecords(forms.map(buildDefaultRecord));
    setMode('register');
  }, [JSON.stringify(forms)]);

  const handleSubmit = (data: any) => {
    console.log(data);
  };

  if (mode === 'register') {
    return (
      <RecipientRegistrationScreenComponent
        forms={forms}
        records={records}
        onSubmit={handleSubmit}
        onCancel={() => {}}
      />
    );
  }

  if (mode === 'review') {
    return null;
  }

  return null;
};
