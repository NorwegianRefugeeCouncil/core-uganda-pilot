import * as React from 'react';
import { VStack, Text, HStack, Link } from 'native-base';
import { FieldKind } from 'core-api-client';

import { NormalisedSubFormFieldValue } from './RecordView.types';

type Props = {
  header: string;
  labels: string[];
  items: NormalisedSubFormFieldValue[][];
};

export const SubformFieldValueComponent: React.FC<Props> = ({
  items,
  labels,
  header,
}) => {
  return (
    <VStack>
      <Text variant="heading" level="3">
        {header}
      </Text>
      <HStack space={3}>
        {labels.map((label, i) => (
          <VStack key={i}>
            <Text variant="body" color="neutral.300">
              {label}
            </Text>
            {items.map((item, j) => {
              const subItem = item[i];
              if (subItem.fieldType === FieldKind.Reference)
                return (
                  <Link key={`${i}-${j}`} href={`/record/${subItem.value}`}>
                    {subItem.formattedValue}
                  </Link>
                );
              return (
                <Text key={`${i}-${j}`} variant="label">
                  {subItem.formattedValue}
                </Text>
              );
            })}
          </VStack>
        ))}
      </HStack>
    </VStack>
  );
};
