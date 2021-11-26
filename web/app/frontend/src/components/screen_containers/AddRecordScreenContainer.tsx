import React from 'react';
import { useForm } from 'react-hook-form';
import { FormDefinition } from 'core-js-api-client/lib/types/types';

import { AddRecordScreen } from '../screens/AddRecordScreen';
import { Platform } from 'react-native';
import { useApiClient } from '../../utils/useApiClient';
import { AddRecordScreenContainerProps } from '../../types/screens';
import {
    getEncryptedLocalData,
    storeEncryptedLocalData,
} from '../../utils/storage';
import { getEncryptionKey } from '../../utils/getEncryptionKey';
import { RECORD_ACTIONS } from '../../reducers/recordsReducers';
import {NetworkStateType} from "expo-network";
import {useNetworkState} from "../../utils/useNetworkState";

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
        client.getForm({ id: formId }).then(data => {
            setForm(data.response);
            setIsLoading(false);
        });
    }, [client]);

    // check for locally stored data on mobile device
    React.useEffect(() => {
        if (!isWeb && recordId) {
            getEncryptedLocalData(recordId).then(data => {
                setHasLocalData(!!data);
                reset(data);
            });
        }
    }, [isWeb, recordId]);

    // react to network changes
    React.useEffect(() => {
        useNetworkState()
            .then(networkState => {
                setIsConnected(networkState.type != NetworkStateType.NONE); // NONE
            })
            .catch(() => setIsLoading(true));
    }, []);

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
