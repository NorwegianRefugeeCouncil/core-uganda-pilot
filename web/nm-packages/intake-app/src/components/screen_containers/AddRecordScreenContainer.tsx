import { FormDefinition } from 'core-js-api-client/lib/types/types';
import { NetworkStateType } from 'expo-network';
import React from 'react';
import { useForm } from 'react-hook-form';
import { Platform } from 'react-native';

import { RECORD_ACTIONS } from '../../reducers/recordsReducers';
import { AddRecordScreenContainerProps } from '../../types/screens';
import { getEncryptionKey } from '../../utils/getEncryptionKey';
import { getNetworkState } from '../../utils/getNetworkState';
import {
    getEncryptedLocalData,
    storeEncryptedLocalData,
} from '../../utils/storage';
import { useApiClient } from '../../utils/useApiClient';
import { AddRecordScreen } from '../screens/AddRecordScreen';

export const AddRecordScreenContainer = ({
    route,
    dispatch,
}: AddRecordScreenContainerProps) => {
    const { formId, recordId } = route.params;

    const client = useApiClient();
    const isWeb = Platform.OS === 'web';

    const [isLoading, setIsLoading] = React.useState(true);
    const [form, setForm] = React.useState<FormDefinition>();
    // const [simulateOffline, setSimulateOffline] = React.useState(!isWeb); // TODO: for testing, remove
    const [isConnected, setIsConnected] = React.useState(false);
    // const [showSnackbar, setShowSnackbar] = React.useState(!isConnected);
    const [hasLocalData, setHasLocalData] = React.useState(false);

    const { control, handleSubmit, formState, reset } = useForm();

    React.useEffect(() => {
        async function fetches() {
            let form, localData, networkState;

            // react to network changes
            try {
                networkState = await getNetworkState();
            } catch (error) {
                console.error(error);
                setIsLoading(true);
            } finally {
                setIsConnected(networkState === NetworkStateType.NONE);
            }

            //
            try {
                const data = await client.getForm({ id: formId });
                form = data?.response;
            } catch (error) {
                console.error(error);
                setIsLoading(true);
            } finally {
                setForm(form);
                setIsLoading(false);
            }

            // check for locally stored data on mobile device
            if (!isWeb && recordId) {
                try {
                    localData = await getEncryptedLocalData(recordId);
                } catch (error) {
                    console.error(error);
                } finally {
                    setHasLocalData(!!localData);
                    reset(localData);
                }
            }
        }
        fetches();
    });

    const onSubmitOffline = async (data: any) => {
        const key = getEncryptionKey();

        storeEncryptedLocalData(recordId, key, data)
            .then(() => {
                setHasLocalData(true);
                dispatch({
                    type: RECORD_ACTIONS.ADD_LOCAL_RECORD,
                    payload: {
                        formId,
                        localRecord: recordId,
                    },
                });
            })
            .catch(() => {
                setHasLocalData(false);
            });
    };

    const onSubmit = (data: any) => {
        handleSubmit(async () => {
            if (isConnected || isWeb) {
                await client.createRecord({ object: { formId, values: data } });
            } else {
                await onSubmitOffline(data);
            }
        });
    };

    return (
        <AddRecordScreen
            form={form}
            control={control}
            onSubmit={onSubmit}
            formState={formState}
            isWeb={isWeb}
            hasLocalData={hasLocalData}
            isConnected={isConnected}
            isLoading={isLoading}
        />
    );
};
