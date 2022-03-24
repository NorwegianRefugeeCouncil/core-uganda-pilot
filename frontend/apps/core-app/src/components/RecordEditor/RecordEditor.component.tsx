import * as React from 'react';
import { Box, VStack, useBreakpointValue } from 'native-base';
import { FieldKind, FormDefinition, getFieldKind } from 'core-api-client';

import { FieldInput } from './FieldInput';

type Props = {
  form: FormDefinition;
  hideKeyFields?: boolean;
};

const useGetFieldWith = () => {
  const width = useBreakpointValue({
    sm: '100%',
    md: '80%',
  });

  return (field: FormDefinition['fields'][number]) => {
    const kind = getFieldKind(field.fieldType);

    if (kind === FieldKind.SubForm) {
      return '100%';
    }

    return width;
  };
};

export const RecordEditorComponent: React.FC<Props> = ({
  form,
  hideKeyFields,
}) => {
  const getFieldWidth = useGetFieldWith();

  return (
    <VStack space={4}>
      {form.fields.map((field) => {
        if (hideKeyFields && field.key) return null;
        return (
          <Box key={field.id} width={getFieldWidth(field)}>
            <FieldInput form={form} field={field} />
          </Box>
        );
      })}
    </VStack>
  );
};
