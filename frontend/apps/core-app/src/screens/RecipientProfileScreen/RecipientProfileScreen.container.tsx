import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { FormWithRecord } from 'core-api-client';

import { RootParamList } from '../../navigation/types';
import { formsClient } from '../../clients/formsClient';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

const DATABASE_ID = '79f2e951-9a42-42d7-9543-e0000231208b'; // Colombia
const FORM_ID = 'ebadab8c-62b3-4191-b304-bd777a091dd3'; // Colombia Individual

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [data, setData] = React.useState<FormWithRecord<Recipient>[]>([]);
  const [prettifiedData, setPrettifiedData] = React.useState<
    FormWithRecord<Recipient>[]
  >([]);
  const [error, setError] = React.useState<string>();

  const prettifyData = (
    originalData: FormWithRecord<Recipient>[],
  ): FormWithRecord<Recipient>[] => {
    // don't show key fields
    const dataWithoutKeys = originalData.map((ancestor) => {
      const noKeys = ancestor.form.fields.filter((field) => !field.key);

      return {
        ...ancestor,
        form: {
          ...ancestor.form,
          fields: noKeys,
        },
      };
    });

    // merge fields of individual into individual beneficiary form
    if (dataWithoutKeys.length >= 1) {
      const mergedIndividual = [
        {
          form: {
            ...dataWithoutKeys[1].form,
            fields: [
              ...dataWithoutKeys[0].form.fields,
              ...dataWithoutKeys[1].form.fields,
            ],
          },
          record: {
            ...dataWithoutKeys[1].record,
            values: [
              ...dataWithoutKeys[0].record.values,
              ...dataWithoutKeys[1].record.values,
            ],
          },
        },
      ];
      return mergedIndividual.concat(dataWithoutKeys.slice(2));
    }
    return dataWithoutKeys;
  };

  React.useEffect(() => {
    (async () => {
      try {
        const recipientData = await formsClient.Recipient.get({
          recordId: route.params.id,
          formId: FORM_ID,
          databaseId: DATABASE_ID,
        });
        setData(recipientData);
      } catch (err) {
        setError(JSON.stringify(err));
      }
      setIsLoading(false);
    })();
  }, [FORM_ID, DATABASE_ID, route.params.id]);

  React.useEffect(() => {
    if (data.length) {
      setPrettifiedData(prettifyData(data));
    }
  }, [data]);

  return (
    <RecipientProfileScreenComponent
      data={prettifiedData}
      isLoading={isLoading}
      error={error}
    />
  );
};
