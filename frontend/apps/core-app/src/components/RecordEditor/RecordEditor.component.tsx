import * as React from 'react';
import { VStack } from 'native-base';
import { FormDefinition, Record } from 'core-api-client';

import { FieldInput } from './FieldInput';

type Props = {
  form: FormDefinition;
  record: Record;
  onChange: () => void;
};

export const RecordEditorComponent: React.FC<Props> = ({
  form,
  record,
  onChange,
}) => {
  return (
    <VStack>
      {form.fields.map((field) => {
        return <FieldInput key={field.id} field={field} />;
      })}
    </VStack>
  );
};
