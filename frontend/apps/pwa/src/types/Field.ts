import { FieldValue } from 'core-api-client';

export type NonSubFormFieldValue = FieldValue & {
  value: Exclude<FieldValue['value'], FieldValue[][]>;
};
