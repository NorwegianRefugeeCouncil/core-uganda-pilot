import * as React from 'react';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Recipient } from 'core-api-client/src/types/client/Recipient';
import { PopulatedForm } from 'core-api-client';

import { RootParamList } from '../../navigation/types';
import { formsClient } from '../../clients/formsClient';

import { RecipientProfileScreenComponent } from './RecipientProfileScreen.component';

const DATABASE_ID = '79f2e951-9a42-42d7-9543-e0000231208b'; // Colombia
const FORM_ID = 'ebadab8c-62b3-4191-b304-bd777a091dd3'; // Colombia Individual

export const RecipientProfileScreenContainer: React.FC = () => {
  const route = useRoute<RouteProp<RootParamList, 'RecipientProfile'>>();
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [data, setData] = React.useState<PopulatedForm<Recipient>[]>([]);
  const [error, setError] = React.useState<string>();

  React.useEffect(() => {
    (async () => {
      try {
        let ancestors = await formsClient.Recipient.getAncestors({
          recordId: route.params.id,
          formId: FORM_ID,
          databaseId: DATABASE_ID,
        });

        // don't show reference fields
        ancestors = ancestors.map((ancestor) => {
          return {
            ...ancestor,
            form: {
              ...ancestor.form,
              fields: ancestor.form.fields.filter(
                (field) => !field.fieldType.reference,
              ),
            },
          };
        });

        // merge fields of individual into individual beneficiary form
        if (ancestors.length >= 1) {
          const mergedIndividual = [
            {
              form: {
                ...ancestors[1].form,
                fields: [
                  ...ancestors[0].form.fields,
                  ...ancestors[1].form.fields,
                ],
              },
              record: {
                ...ancestors[1].record,
                values: [
                  ...ancestors[0].record.values,
                  ...ancestors[1].record.values,
                ],
              },
            },
          ];
          setData(mergedIndividual.concat(ancestors.slice(2)));
        } else {
          setData(ancestors);
        }
      } catch (err) {
        setError(JSON.stringify(err));
      }
      setIsLoading(false);
    })();
  }, [FORM_ID, DATABASE_ID, route.params.id]);

  return (
    <RecipientProfileScreenComponent
      data={data}
      isLoading={isLoading}
      error={error}
    />
  );
};
