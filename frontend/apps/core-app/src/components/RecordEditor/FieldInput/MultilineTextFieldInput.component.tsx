import * as React from 'react';
import { FormControl, TextArea } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition } from 'core-api-client';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const MultilineTextFieldInput: React.FC<Props> = ({ formId, field }) => {
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
      <FormControl.Label size="xs">{field.name}</FormControl.Label>
      <TextArea
        ref={ref}
        onBlur={onBlur}
        onChangeText={onChange}
        value={value}
      />
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error}</FormControl.ErrorMessage>
    </FormControl>
  );
};
