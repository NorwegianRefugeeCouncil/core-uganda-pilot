import React from 'react';
import {common, layout} from '../../styles';
import {Button, Platform, ScrollView, Text, View} from 'react-native';
import {useForm} from "react-hook-form";
import {Individual} from "../../../../client/src/types/models";
import {Subject} from 'rxjs';
import iamClient from "../../utils/clients";
import {PartyAttributeDefinition, PartyAttributeDefinitionList} from "core-js-api-client/lib/types/models";
import _ from 'lodash';
import {createIndividual} from "../../services/individuals";
import FormControl from "../form/FormControl";
import * as SecureStore from 'expo-secure-store';
import * as Network from 'expo-network';
import {Snackbar, Switch} from 'react-native-paper';

export interface FlatIndividual extends Omit<Individual, 'attributes'> {
    attributes: { [p: string]: string }
}

const IndividualScreen: React.FC<any> = ({route}) => {
    const {id} = route.params;
    const isWeb = Platform.OS === 'web';

    const [isLoading, setIsLoading] = React.useState(true);
    const [simulateOffline, setSimulateOffline] = React.useState(!isWeb); // TODO: for testing, remove
    const [attributes, setAttributes] = React.useState<PartyAttributeDefinition[]>([]);
    const [individual, setIndividual] = React.useState<FlatIndividual>();
    const [isConnected, setIsConnected] = React.useState(!simulateOffline);
    const [showSnackbar, setShowSnackbar] = React.useState(!isConnected);
    const [hasLocalData, setHasLocalData] = React.useState(false);

    let attributesSubject = new Subject([]);
    let individualsSubject = new Subject([]);
    let flatAttributes: { [p: string]: string } = {};
    let submitData: { [p: string]: string[] } = {};

    const {control, handleSubmit, formState, reset} = useForm({
        defaultValues: individual
    });

    const onSubmitOffline = async (data: string[]) => {
        SecureStore.setItemAsync(id, JSON.stringify(data))
            .then((data) => {
                console.log('SBUMITTED!!!', data)
                setHasLocalData(true)
            })
            .catch((e) => {
                console.log('SUBMIT ERROR', e)
                setHasLocalData(false)
            });

    }
    const onSubmit = (data: string[]) => {

        // wrap attributes in arrays, move somewhere else
        _(data).forEach((value, key) => {
            submitData[key] = [value];
        });

        if (isConnected || isWeb) {
            createIndividual({
                id,
                attributes: submitData,
                partyTypeIds: individual?.partyTypeIds || []
            });
        } else {
            onSubmitOffline(data);
        }

    };

    {/* check for locally stored data on mobile device */    }
    React.useEffect(() => {
        if (!isWeb) {
            SecureStore.getItemAsync(id)
                .then((data) => {
                    setHasLocalData(data != null)
                    const newIndividual: FlatIndividual = {
                        id: id,
                        partyTypeIds: individual?.partyTypeIds || [],
                        attributes: data == null ? individual?.attributes : JSON.parse(data)
                    };
                    reset(newIndividual);
                    // TODO: delete data, once extracted to save space. or only after submit?
                })
        }
    }, [isWeb])

    // react to network changes
    React.useEffect(() => {
        Network.getNetworkStateAsync()
            .then((networkState) => {
                // TODO: uncomment, use real network state
                // setIsConnected(networkState.type != NetworkStateType.NONE); // NONE
            })
            .catch(() => setIsLoading(true))
    }, [simulateOffline])

    // get data for form and individual
    React.useEffect(() => {
        attributesSubject.pipe(iamClient.PartyAttributeDefinitions().List()).subscribe(
            (data: PartyAttributeDefinitionList) => {
                setAttributes(data.items)
            }
        );
        individualsSubject.pipe(iamClient.Individuals().Get()).subscribe(
            (data: Individual) => {
                _(data.attributes).forEach((value, key) => {
                    flatAttributes[key] = value[0];
                });
                const flatIndividual = {...data, attributes: flatAttributes};
                setIndividual(flatIndividual)
            }
        );
        attributesSubject.next(null);
        individualsSubject.next(id);

        return () => {
            if (attributesSubject) {
                attributesSubject.unsubscribe();
            }
            if (individualsSubject) {
                individualsSubject.unsubscribe();
            }
        };
    }, []);

    // check if data has been received
    React.useEffect(() => {
        if (individual) {
            setIsLoading(false);
        }
    }, [individual])

    return (
        <View style={{paddingBottom: 26}}>
            {/* simulate network changes, for testing */}
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
            <ScrollView>
                {!isLoading && (
                    <View style={[layout.container, layout.body, common.darkBackground]}>
                        {attributes.map((a) =>
                            <FormControl
                                key={a.id}
                                formControl={a.formControl}
                                style={{width: '100%'}}
                                value={individual?.attributes[a.id]}
                                control={control}
                                name={`attributes.${a.id}`}
                                errors={formState.errors}
                            />
                        )}
                        <Button
                            title="Submit"
                            onPress={handleSubmit(onSubmit)}
                        />

                    </View>
                )}
            </ScrollView>
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
        </View>
    );
};

export default IndividualScreen;
