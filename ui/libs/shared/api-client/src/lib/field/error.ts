import { Path } from './index';

export type ErrorType =
  'FieldValueNotFound'
  | 'FieldValueRequired'
  | 'FieldValueDuplicate'
  | 'FieldValueInvalid'
  | 'FieldValueNotSupported'
  | 'FieldValueForbidden'
  | 'FieldValueTooLong'
  | 'FieldValueTooMany'
  | 'InternalError'


export class Error {
  public constructor(
    public type: ErrorType,
    public field: string,
    public badValue: any,
    public detail: string
  ) {
  }
}

export type ErrorList = Error[]

export const NotFound = (field: Path, value: any): Error => {
  return new Error('FieldValueNotFound', field.string(), value, '');
};

export const Required = (field: Path, detail: string): Error => {
  return new Error('FieldValueRequired', field.string(), '', detail);
};

export const Duplicate = (field: Path, value: any): Error => {
  return new Error('FieldValueDuplicate', field.string(), value, '');
};

export const Invalid = (field: Path, value: any, detail: string) => {
  return new Error('FieldValueInvalid', field.string(), value, detail);
};

