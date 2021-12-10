import React from 'react';

import { RECORD_ACTIONS } from '../../reducers/recordsReducers';
import { RecordsScreenContainerProps } from '../../types/screens';
import RecordsScreen from '../screens/RecordsScreen';
import useApiClient from '../../utils/clients';

export const RecordsScreenContainer = ({ navigation, route, state, dispatch }: RecordsScreenContainerProps) => {
  const { formId, databaseId } = route.params;
  const [isLoading, setIsLoading] = React.useState(true);
  const apiClient = useApiClient();

  React.useEffect(() => {
    if (formId && databaseId) {
      apiClient
        .listRecords({ formId, databaseId })
        .then((data) => {
          dispatch({
            type: RECORD_ACTIONS.GET_RECORDS,
            payload: {
              formId,
              records: data.response?.items,
            },
          });
        })
        .catch(console.error)
        .finally(() => setIsLoading(false));
    }
  }, [formId, databaseId]);

  const recordsScreenProps = { isLoading, state, navigation, route };
  return <RecordsScreen {...recordsScreenProps} />;
};
