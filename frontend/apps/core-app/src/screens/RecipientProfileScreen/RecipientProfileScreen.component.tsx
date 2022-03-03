import * as React from 'react';
import { Text } from 'native-base';
import { RouteProp } from '@react-navigation/native';
import { FormDefinition, Record } from 'core-api-client';

import { RootParamList } from '../../navigation/types';
import { RecordView } from '../../components/RecordView';
import { formsClient } from '../../clients/formsClient';

import * as Styles from './RecipientProfileScreen.styles';

type Props = {
  route: RouteProp<RootParamList, 'RecipientProfile'>;
};

// Change the ids to something in your local database
const DATABASE_ID = '529522ff-5e09-4bb9-a703-558548489b93';
const FORM_ID = '74b90de4-1e23-467c-affc-e41657cb67cb';
const RECORD_ID = '717162d8-ca71-47c4-aa19-c044a37c3dd7';

const useFakeRecord = (
  formId: string,
  databaseId: string,
  recordId: string,
): [FormDefinition | null, Record | null] => {
  const [form, setForm] = React.useState<FormDefinition | null>(null);
  const [record, setRecord] = React.useState<Record | null>(null);

  React.useEffect(() => {
    (async () => {
      const formResponse = await formsClient.Form.get({ id: formId });
      if (formResponse.response) setForm(formResponse.response);

      const recordListResponse = await formsClient.Record.get({
        recordId,
        formId,
        databaseId,
      });
      if (recordListResponse.response) setRecord(recordListResponse.response);
    })();
  }, [formId, databaseId, recordId]);

  return [form, record];
};

export const RecipientProfileScreenComponent: React.FC<Props> = ({ route }) => {
  const [form, record] = useFakeRecord(FORM_ID, DATABASE_ID, RECORD_ID);

  return (
    <Styles.Container>
      <Text variant="display">
        {route.name}: {route.params.id}
      </Text>
      {form && record && <RecordView form={form} record={record} />}
    </Styles.Container>
  );
};
