import * as React from 'react';
import { VStack } from 'native-base';
import { FormWithRecord } from 'core-api-client';
import { Accordion } from 'core-design-system';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { RecordEditor } from '../../RecordEditor';

type Props = {
  data: FormWithRecord<Recipient>[];
};

export const RecipientEditorComponent: React.FC<Props> = ({ data }) => (
  <VStack space={4}>
    <Accordion header={data[1].form.name} defaultOpen>
      <VStack space={4}>
        <RecordEditor form={data[0].form} hideKeyFields />
        <RecordEditor form={data[1].form} hideKeyFields />
      </VStack>
    </Accordion>
    {data.slice(2).map(({ form }) => (
      <Accordion key={form.id} header={form.name} defaultOpen>
        <RecordEditor form={form} hideKeyFields />
      </Accordion>
    ))}
  </VStack>
);
