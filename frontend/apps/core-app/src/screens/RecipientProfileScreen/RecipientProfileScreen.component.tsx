import * as React from 'react';
import { Skeleton, Text } from 'native-base';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormDefinition } from 'core-api-client';

import { RecordView } from '../../components/RecordView';

import * as Styles from './RecipientProfileScreen.styles';

type Props = {
  recipient: Recipient | null;
  isLoading: boolean;
  form: FormDefinition | null;
  error: any;
};

export const RecipientProfileScreenComponent: React.FC<Props> = ({
  recipient,
  isLoading,
  form,
  error,
}) => {
  return (
    <Styles.Container>
      {isLoading && <Skeleton h="20" p="4" />}
      {error && (
        <Text variant="heading" color="signalDanger">
          {JSON.stringify(error)}
        </Text>
      )}
      {form && recipient && <RecordView form={form} record={recipient} />}
    </Styles.Container>
  );
};
