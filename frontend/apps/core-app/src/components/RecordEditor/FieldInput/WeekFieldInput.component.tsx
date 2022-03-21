import * as React from 'react';
import { FormControl, Text } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition } from 'core-api-client';
import { Platform } from 'expo-modules-core';

import { Input } from '../../Web/Input';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const WeekFieldInput: React.FC<Props> = ({ formId, field }) => {
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
    <FormControl isInvalid={invalid}>
      <FormControl.Label>{field.name}</FormControl.Label>
      {Platform.OS === 'web' ? (
        <Input
          type="week"
          value={value}
          onChange={onChange}
          invalid={invalid}
        />
      ) : (
        <Text>Not implemented on mobile</Text>
      )}
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
    </FormControl>
  );
};
