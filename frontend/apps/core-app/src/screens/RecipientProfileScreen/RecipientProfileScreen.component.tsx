import * as React from 'react';
import { ScrollView, Skeleton, Text } from 'native-base';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormWithRecord } from 'core-api-client';

import { RecipientViewer } from '../../components/Recipient/RecipientViewer';

import * as Styles from './RecipientProfileScreen.styles';

type Props = {
  isLoading: boolean;
  data: FormWithRecord<Recipient>[];
  error: string | null;
};

export const RecipientProfileScreenComponent: React.FC<Props> = ({
  data,
  isLoading,
  error,
}) => (
  <ScrollView width="100%">
    <Styles.Container width="100%" maxWidth="1180px" marginX="auto">
      {isLoading && !error && <Skeleton h="20" p="4" />}
      {error && (
        <Text variant="body" color="signalDanger">
          {error}
        </Text>
      )}
      {!isLoading && !error && <RecipientViewer data={data} />}
    </Styles.Container>
  </ScrollView>
);
