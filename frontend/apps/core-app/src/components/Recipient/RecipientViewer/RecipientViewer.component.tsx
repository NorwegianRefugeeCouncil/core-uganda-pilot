import * as React from 'react';
import { VStack } from 'native-base';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormWithRecord } from 'core-api-client';
import { Accordion } from 'core-design-system';

import { RecordView } from '../../RecordView';

type Props = {
  data: FormWithRecord<Recipient>[];
};

export const RecipientViewerComponent: React.FC<Props> = ({ data }) => (
  <VStack space={4}>
    <Accordion header={data[1].form.name} defaultOpen>
      <VStack space={4}>
        <RecordView form={data[0].form} record={data[0].record} hideKeyFields />
        <RecordView form={data[1].form} record={data[1].record} hideKeyFields />
      </VStack>
    </Accordion>
    {data.slice(2).map(({ form, record: recipient }) => (
      <Accordion key={form.id} header={form.name} defaultOpen>
        <RecordView form={form} record={recipient} hideKeyFields />
      </Accordion>
    ))}
  </VStack>
);
