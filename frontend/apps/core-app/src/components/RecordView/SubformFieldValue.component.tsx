import * as React from 'react';
import { VStack, Text } from 'native-base';
import { Column } from 'react-table';

import { SubFormTable } from '../SubFormTable';

type Props = {
  header: string;
  data: Record<string, string>[];
  columns: Column<Record<string, string>>[];
};

export const SubformFieldValueComponent: React.FC<Props> = ({
  header,
  data,
  columns,
}) => {
  return (
    <VStack>
      <Text variant="heading" level="3">
        {header}
      </Text>
      <SubFormTable data={data} columns={columns} />
    </VStack>
  );
};
