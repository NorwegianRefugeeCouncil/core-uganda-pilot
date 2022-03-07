import { FieldDefinition, FieldValue } from 'core-api-client';
import { FieldErrors } from 'react-hook-form';

import { FormValue } from '../../reducers/Recorder/types';

export type FieldEditorProps = {
  field: FieldDefinition;
  value: FieldValue['value'];
  onChange: (value: string | string[] | null) => void;
  onAddSubRecord: () => void;
  onSelectSubRecord: (subRecordId: string) => void;
  subRecords: FormValue[] | undefined;
  errors: FieldErrors;
};
