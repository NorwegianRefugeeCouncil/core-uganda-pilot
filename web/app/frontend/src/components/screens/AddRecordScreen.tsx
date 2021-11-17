import React from 'react';
import {Snackbar, Switch} from 'react-native-paper';
import {common, layout} from '../../styles';
import {Button, Platform, ScrollView, Text, View} from 'react-native';
import useApiClient from "../../utils/clients";
import {FormDefinition} from "core-js-api-client/lib/types/types";
import {useForm} from "react-hook-form";
import * as Network from "expo-network";
import FormControl from "../form/FormControl";
import {getEncryptedLocalData, storeEncryptedLocalData} from "../../utils/storage";
import {RECORD_ACTIONS} from "../../reducers/recordsReducers";
import {getEncryptionKey} from "../../utils/getEncryptionKey";

const AddRecordScreen: React.FC<any> = ({route, dispatch}) => {
    const isWeb = Platform.OS === 'web';
    const {formId, recordId} = route.params;

    const [isLoading, setIsLoading] = React.useState(true);
    const [form, setForm] = React.useState<FormDefinition>();
    const [simulateOffline, setSimulateOffline] = React.useState(!isWeb); // TODO: for testing, remove
    const [isConnected, setIsConnected] = React.useState(!simulateOffline);
    const [showSnackbar, setShowSnackbar] = React.useState(!isConnected);
    const [hasLocalData, setHasLocalData] = React.useState(false);

    const client = useApiClient();
    const {control, handleSubmit, formState, reset} = useForm();

    React.useEffect(() => {
        client.getForm({id: formId})
            .then((data) => {
                setForm(data.response)
                setIsLoading(false)
            })
        // TODO add catch
    }, []);

    const onSubmitOffline = async (data: any) => {
        const key = getEncryptionKey();

        storeEncryptedLocalData(recordId, key, data)
            .then(() => {
                setHasLocalData(true)
                dispatch({
                    type: RECORD_ACTIONS.ADD_LOCAL_RECORD, payload: {
                        formId,
                        localRecord: recordId
                    }
                })
            })
            .catch(() => {
                setHasLocalData(false)
            });

    }
    const onSubmit = (data: any) => {
        if (isConnected || isWeb) {
            client.createRecord({object: {formId, values: data}})
        } else {
            onSubmitOffline(data);
        }
    };

    // check for locally stored data on mobile device
    React.useEffect(() => {
        if (!isWeb && recordId) {

            getEncryptedLocalData(recordId)
                .then((data) => {
                    setHasLocalData(!!data);
                    reset(data);
                });
            }
    }, [isWeb, recordId])

    // react to network changes
    React.useEffect(() => {
        Network.getNetworkStateAsync()
            .then((networkState) => {
                // TODO: uncomment, use real network state
                // setIsConnected(networkState.type != NetworkStateType.NONE); // NONE
            })
            .catch(() => setIsLoading(true))
    }, [simulateOffline])

    return (
        <ScrollView contentContainerStyle={[layout.container, layout.body, common.darkBackground]}>

            <View style={[]}>
                {/* simulate network changes, for testing */}
                {!isWeb && (
                    <View style={{display: "flex", flexDirection: "row"}}>
                        <Switch
                            value={simulateOffline}
                            onValueChange={() => {
                                setSimulateOffline(!simulateOffline)
                                setIsConnected(simulateOffline)
                                setShowSnackbar(!simulateOffline)
                            }}
                        />
                        <Text> simulate being offline </Text>
                    </View>
                )}

                {/* upload data collected offline */}
                {hasLocalData && (
                    <View style={{display: "flex", flexDirection: "column"}}>
                        <Text>
                            There is locally stored data for this individual.
                        </Text>
                    </View>
                )}
                {hasLocalData && isConnected && (
                    <View style={{display: "flex", flexDirection: "column"}}>
                        <Text>
                            Do you want to upload it?
                        </Text>
                        <Button
                            title="Submit local data"
                            onPress={handleSubmit(onSubmit)}
                        />
                    </View>
                )}
                {!isLoading && (
                    <View style={{width: '100%'}}>
                        {form?.fields.map((field) => {
                            return (
                                <FormControl
                                    key={field.code}
                                    fieldDefinition={field}
                                    style={{width: '100%'}}
                                    // value={''} // take value from record
                                    control={control}
                                    name={field.id}
                                    errors={formState.errors}
                                />
                            )
                        })}
                        <Button
                            title="Submit"
                            onPress={handleSubmit(onSubmit)}
                        />
                    </View>
                )}
            </View>
            <Snackbar
                visible={showSnackbar}
                onDismiss={() => setShowSnackbar(false)}
                action={{
                    label: 'Got it',
                    onPress: () => setShowSnackbar(false)
                }}
            >
                No internet connection. Submitted data will be stored locally.
            </Snackbar>
        </ScrollView>
    );
};

export default AddRecordScreen;
