import { FieldDefinition } from 'core-api-client';

import { FormValue } from '../../reducers/recorder';

export type FieldEditorProps = {
  field: FieldDefinition;
  value: string | string[] | null;
  setValue: (value: string | string[] | null) => void;
  addSubRecord: () => void;
  selectSubRecord: (subRecordId: string) => void;
  subRecords: FormValue[] | undefined;
};
