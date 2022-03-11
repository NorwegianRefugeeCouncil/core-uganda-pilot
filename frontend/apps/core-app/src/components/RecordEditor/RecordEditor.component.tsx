import * as React from 'react';
import { Box, Stack, useBreakpointValue } from 'native-base';
import { FieldKind, FormDefinition, getFieldKind } from 'core-api-client';

import { FieldInput } from './FieldInput';

type Props = {
  form: FormDefinition;
  direction: 'row' | 'column';
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
  direction = 'column',
}) => {
  const getFieldWidth = useGetFieldWith();

  return (
    <Stack direction={direction} space={4} flexWrap="wrap">
      {form.fields.map((field) => {
        return (
          <Box key={field.id} width={getFieldWidth(field)}>
            <FieldInput form={form} field={field} />
          </Box>
        );
      })}
    </Stack>
  );
};
