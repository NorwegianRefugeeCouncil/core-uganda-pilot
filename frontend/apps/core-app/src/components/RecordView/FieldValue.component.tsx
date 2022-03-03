import * as React from 'react';
import { Stack, Text, Link, useBreakpointValue } from 'native-base';
import { FieldKind } from 'core-api-client';

import { NormalisedBasicField } from './RecordView.types';

type Props = {
  item: NormalisedBasicField;
};

const useResponsiveStyles = () => {
  const direction = useBreakpointValue({
    sm: 'column',
    md: 'row',
  });

  const alignItems = useBreakpointValue({
    sm: 'flex-start',
    md: 'flex-start',
  });

  const justifyContent = useBreakpointValue({
    sm: 'center',
    md: 'flex-start',
  });

  const space = useBreakpointValue({
    sm: '1',
    md: '2',
  });

  const labelWidth = useBreakpointValue({
    sm: '100%',
    md: '40%',
  });

  const valueWidth = useBreakpointValue({
    sm: '100%',
    md: '60%',
  });

  return {
    direction,
    alignItems,
    justifyContent,
    space,
    labelWidth,
    valueWidth,
  };
};

export const FieldValueComponent: React.FC<Props> = ({ item }) => {
  const styles = useResponsiveStyles();

  return (
    <Stack
      direction={styles.direction}
      alignItems={styles.alignItems}
      justifyContent={styles.justifyContent}
      space={styles.space}
      width="100%"
    >
      <Text variant="body" color="neutral.300" width={styles.labelWidth}>
        {item.label}
      </Text>
      {item.fieldType === FieldKind.Reference ? (
        <Link href={`/record/${item.value}`} width={styles.valueWidth}>
          {item.formattedValue}
        </Link>
      ) : (
        <Text variant="label" width="60%">
          {item.formattedValue}
        </Text>
      )}
    </Stack>
  );
};
