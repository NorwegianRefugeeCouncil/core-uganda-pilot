import * as React from 'react';
import { FormControl, Checkbox } from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition, Validation } from 'core-api-client';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};
export const CheckboxFieldInput: React.FC<Props> = ({ formId, field }) => {
  const { control } = useFormContext();

  const {
    field: { onChange, value },
    fieldState: { error, invalid },
  } = useController({
    name: `${formId}.${field.id}`,
    control,
    rules: Validation.Record.formValidationRules.field.checkbox(field),
  });

  const handleChange = (checked: boolean) => {
    onChange(checked ? 'true' : 'false');
  };

  return (
    <FormControl isRequired={field.required} isInvalid={invalid}>
      <Checkbox
        testID="checkbox-field-input"
        value={field.id}
        onChange={handleChange}
        isChecked={value === 'true'}
      >
        {field.name}
      </Checkbox>
      <FormControl.HelperText>{field.description}</FormControl.HelperText>
      <FormControl.ErrorMessage>{error?.message}</FormControl.ErrorMessage>
    </FormControl>
  );
};
