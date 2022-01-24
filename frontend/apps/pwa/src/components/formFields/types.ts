import { FieldDefinition } from 'core-api-client';

import { FormValue } from '../../reducers/Recorder/types';

export type FieldEditorProps = {
  field: FieldDefinition;
  value: string | string[] | null;
  onChange: (value: string | string[] | null) => void;
  onAddSubRecord: () => void;
  onSelectSubRecord: (subRecordId: string) => void;
  subRecords: FormValue[] | undefined;
};
