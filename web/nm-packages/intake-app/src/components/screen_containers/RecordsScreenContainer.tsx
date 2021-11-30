import React from 'react';
import { RecordsScreenContainerProps } from '../../types/screens';
import {RecordsScreen} from '../screens/RecordsScreen';
import { useApiClient } from '../../utils/useApiClient';
import { RECORD_ACTIONS } from '../../reducers/recordsReducers';

export const RecordsScreenContainer = ({
    navigation,
    route,
    state,
    dispatch,
}: RecordsScreenContainerProps) => {
    const client = useApiClient();
    const { formId, databaseId } = route.params;
    const [isLoading, setIsLoading] = React.useState(true);
    const { records, localRecords } = state.formsById[formId];

    React.useEffect(() => {
        client.listRecords({ formId, databaseId }).then(data => {
            dispatch({
                type: RECORD_ACTIONS.GET_RECORDS,
                payload: {
                    formId,
                    records: data.response?.items,
                },
            });
            setIsLoading(false);
        });
    }, []);
    return (
        <RecordsScreen
            navigation={navigation}
            isLoading={isLoading}
            records={records}
            localRecords={localRecords}
            formId={formId}
        />
    );
};
