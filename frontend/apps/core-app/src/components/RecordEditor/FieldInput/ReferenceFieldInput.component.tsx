import * as React from 'react';
import { FormControl, Select } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition, Record } from 'core-api-client';

import { formsClient } from '../../../clients/formsClient';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const ReferenceFieldInput: React.FC<Props> = ({ formId, field }) => {
  const [referenceRecords, setReferenceRecords] = React.useState<Record[]>([]);

  React.useEffect(() => {
    (async () => {
      try {
        if (field.fieldType.reference) {
          const r = await formsClient.Record.list({
            formId: field.fieldType.reference.formId,
            databaseId: field.fieldType.reference.databaseId,
          });
          if (r.response) setReferenceRecords(r.response.items);
          else setReferenceRecords([]);
        } else {
          setReferenceRecords([]);
        }
      } catch {
        setReferenceRecords([]);
      }
    })();
  }, [field.fieldType.reference?.formId]);

  const { control } = useFormContext();

  const {
    field: { onChange, value },
    fieldState: { error, invalid },
  } = useController({
    name: `${formId}.${field.id}`,
    control,
    rules: {}, // TODO Record validation
  });

  return (
    <FormControl isRequired={field.required} isInvalid={invalid}>
      <FormControl.Label>{field.name}</FormControl.Label>
      <Select
        testID="reference-field-input"
        mt="1"
        onValueChange={onChange}
        selectedValue={value}
      >
        {referenceRecords.map((r, i) => (
          <Select.Item
            testID={`reference-field-input-option-${i}`}
            key={r.id}
            label={r.id}
            value={r.id}
          />
        ))}
      </Select>
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
    </FormControl>
  );
};
