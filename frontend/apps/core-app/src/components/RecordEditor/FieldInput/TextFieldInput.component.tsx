import * as React from 'react';
import { FormControl, Input } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition } from 'core-api-client';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const TextFieldInput: React.FC<Props> = ({ formId, field }) => {
  const { control } = useFormContext();

  const {
    field: { onChange, onBlur, value, ref },
    fieldState: { error, invalid },
  } = useController({
    name: `${formId}.${field.id}`,
    control,
    rules: {}, // TODO Record validation
  });

  return (
    <FormControl isInvalid={invalid}>
      <FormControl.Label>{field.name}</FormControl.Label>
      <Input
        testID="text-field-input"
        ref={ref}
        onBlur={onBlur}
        onChangeText={onChange}
        value={value}
        autoCompleteType="off"
      />
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
    </FormControl>
  );
};
