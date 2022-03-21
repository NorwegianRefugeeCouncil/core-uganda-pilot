import * as React from 'react';
import { FormControl, Select } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition } from 'core-api-client';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const SingleSelectFieldInput: React.FC<Props> = ({ formId, field }) => {
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
        testID="single-select-field-input"
        mt="1"
        onValueChange={onChange}
        selectedValue={value}
      >
        {field.fieldType.singleSelect?.options.map((option, i) => (
          <Select.Item
            key={option.id}
            testID={`single-select-field-input-option-${i}`}
            label={option.name}
            value={option.id}
          />
        ))}
      </Select>
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
    </FormControl>
  );
};
