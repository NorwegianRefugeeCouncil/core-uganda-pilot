import { FormSchema } from './form-types';

const schema: FormSchema = {
  kind: 'IntakeForm',
  apiVersion: 'core.nrc.no/v1',
  metadata: {
    name: 'Beneficiary Details',
    uid: '9284012983120983asojdajdaoj',
    creationTimestamp: '2020-05-01',
  },
  data: {
    firstName: {
      type: 'shortText',
      children: [],
      options: {
        name: {
          en: 'First Name',
        },
        description: {
          en: "The beneficiary's first name",
        },
        tooltip: {
          en:
            'Enter the name the beneficiary prefers to use as their first name',
        },
      },
    },
    lastName: {
      type: 'shortText',
      children: [],
      options: {
        name: {
          en: 'Last Name',
        },
        description: {
          en: "The beneficiary's last name",
        },
        tooltip: {
          en:
            'Enter the name the beneficiary prefers to use as their last name',
        },
      },
    },
    consent: {
      type: 'checkbox',
      children: [],
      options: {
        name: {
          en: 'demo checkbox field',
        },
        description: {
          en: 'demo checkbox field description',
        },
        tooltip: {
          en: 'this is a demo checkbox field',
        },
      },
    },
  },
};

export async function getSchema() {
  return new Promise((resolve) => resolve(schema));
}
