export type FormElementType =
  | 'section'
  | 'shortText'
  | 'longText'
  | 'integer'
  | 'select'
  | 'date'
  | 'dateTime'
  | 'time'
  | 'checkbox';

export interface i18nField {
  [locale: string]: string;
}

export interface FormElementOptions {
  name: i18nField;
  description: i18nField;
  tooltip: i18nField;
}

export interface FormElement {
  type: FormElementType;
  children?: Array<FormElement>;
  options: FormElementOptions;
}

export interface FormDataField {
  [field: string]: FormElement;
}

export interface Metadata {
  name: string;
  uid: string;
  creationTimestamp: string;
}

export interface FormSchema {
  kind: string;
  apiVersion: string;
  metadata: Metadata;
  data: FormDataField;
}
