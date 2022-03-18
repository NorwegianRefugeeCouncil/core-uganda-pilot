import { FieldKind } from 'core-api-client';

export type NormalisedBasicField = {
  key: boolean;
  fieldType: Exclude<FieldKind, FieldKind.SubForm>;
  value: string | string[] | null;
  formattedValue: string;
  label: string;
};

export type NormalisedSubFormFieldValue = Omit<NormalisedBasicField, 'label'>;

export type NormalisedSubFormField = {
  key: boolean;
  fieldType: FieldKind.SubForm;
  header: string;
  labels: string[];
  values: NormalisedSubFormFieldValue[][];
};

export type NormalisedFieldValue =
  | NormalisedSubFormField
  | NormalisedBasicField;
