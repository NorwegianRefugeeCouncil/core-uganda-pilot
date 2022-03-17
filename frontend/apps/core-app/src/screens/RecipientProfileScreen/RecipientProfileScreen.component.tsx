import * as React from 'react';
import { ScrollView, Skeleton, Text } from 'native-base';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormWithRecord } from 'core-api-client';
import { Accordion } from 'core-design-system';

import { RecordView } from '../../components/RecordView';

import * as Styles from './RecipientProfileScreen.styles';

type Props = {
  isLoading: boolean;
  data: FormWithRecord<Recipient>[];
  error?: string;
};

export const RecipientProfileScreenComponent: React.FC<Props> = ({
  data,
  isLoading,
  error,
}) => {
  return (
    <ScrollView>
      <Styles.Container>
        {isLoading && <Skeleton h="20" p="4" />}
        {error && (
          <Text variant="heading" color="signalDanger">
            {error}
          </Text>
        )}
        {data.map(({ form, record: recipient }) => {
          return (
            <Accordion header={form.name || ''} key={form.id}>
              <RecordView form={form} record={recipient} />
            </Accordion>
          );
        })}
      </Styles.Container>
    </ScrollView>
  );
};
