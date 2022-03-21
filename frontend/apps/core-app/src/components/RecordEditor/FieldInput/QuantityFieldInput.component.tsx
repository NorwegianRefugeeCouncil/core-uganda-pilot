import * as React from 'react';
import { FormControl, Input } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition } from 'core-api-client';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const QuantityFieldInput: React.FC<Props> = ({ formId, field }) => {
  const { control } = useFormContext();

  const {
    field: { onChange, onBlur, value, ref },
    fieldState: { error, invalid },
  } = useController({
    name: `${formId}.${field.id}`,
    control,
    rules: {}, // TODO Record validation
  });

  const handleOnChange = (v: string) => {
    // if value is a number or a decimal
    if (v.match(/^[0-9]*\.?[0-9]*$/)) {
      onChange(v);
    }
  };

  return (
    <FormControl isInvalid={invalid}>
      <FormControl.Label>{field.name}</FormControl.Label>
      <Input
        testID="quantity-field-input"
        ref={ref}
        onBlur={onBlur}
        onChangeText={handleOnChange}
        value={value}
        autoCompleteType="off"
      />
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
    </FormControl>
  );
};
