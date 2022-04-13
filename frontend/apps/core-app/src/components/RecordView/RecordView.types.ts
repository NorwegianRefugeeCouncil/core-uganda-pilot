import { FieldKind } from 'core-api-client';
import { Column } from 'react-table';

export type NormalisedBasicField = {
  key: boolean;
  fieldType: Exclude<FieldKind, FieldKind.SubForm>;
  value: string | string[] | null;
  formattedValue: string;
  label: string;
};

export type NormalisedSubFormField = {
  key: boolean;
  fieldType: FieldKind.SubForm;
  header: string;
  data: Record<string, string>[];
  columns: Column<Record<string, string>>[];
};

export type NormalisedFieldValue =
  | NormalisedSubFormField
  | NormalisedBasicField;
