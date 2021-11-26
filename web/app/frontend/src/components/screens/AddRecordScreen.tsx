import React from 'react';
import { Button, ScrollView, Text, View } from 'react-native';
import { common, layout } from '../../styles';
import FormControl from '../form/FormControl';
import { FormDefinition } from 'core-js-api-client/lib/types/types';
import { Control, FieldValues, FormState } from 'react-hook-form';

export type AddRecordScreenProps = {
    form?: FormDefinition;
    control: Control;
    onSubmit: (data: any) => void;
    formState: FormState<FieldValues>;
    isWeb: boolean;
    hasLocalData: boolean;
    isConnected: boolean;
    isLoading: boolean;
};

export const AddRecordScreen = ({
    form,
    control,
    onSubmit,
    formState,
    isWeb,
    hasLocalData,
    isConnected,
    isLoading,
}: AddRecordScreenProps) => {
    return (
        <ScrollView
            contentContainerStyle={[
                layout.container,
                layout.body,
                common.darkBackground,
            ]}
        >
            <View>
                {/* simulate network changes, for testing */}
                {!isWeb && (
                    <View style={{ display: 'flex', flexDirection: 'row' }}>
                        {/*<Switch*/}
                        {/*    value={simulateOffline}*/}
                        {/*    onValueChange={() => {*/}
                        {/*        setSimulateOffline(!simulateOffline);*/}
                        {/*        setIsConnected(simulateOffline);*/}
                        {/*        setShowSnackbar(!simulateOffline);*/}
                        {/*    }}*/}
                        {/*/>*/}
                        <Text> simulate being offline </Text>
                    </View>
                )}

                {/* upload data collected offline */}
                {hasLocalData && (
                    <View style={{ display: 'flex', flexDirection: 'column' }}>
                        <Text>
                            There is locally stored data for this individual.
                        </Text>
                    </View>
                )}
                {hasLocalData && isConnected && (
                    <View style={{ display: 'flex', flexDirection: 'column' }}>
                        <Text>Do you want to upload it?</Text>
                        <Button title="Submit local data" onPress={onSubmit} />
                    </View>
                )}
                {!isLoading && (
                    <View style={{ width: '100%' }}>
                        <Text>{form?.name}</Text>
                        {form?.fields.map(field => {
                            return (
                                <FormControl
                                    key={field.code}
                                    fieldDefinition={field}
                                    style={{ width: '100%' }}
                                    // value={''} // take value from record
                                    control={control}
                                    name={field.id}
                                    errors={formState.errors}
                                />
                            );
                        })}
                        <Button title="Submit" onPress={onSubmit} />
                    </View>
                )}
            </View>
            {/*<Snackbar*/}
            {/*    visible={showSnackbar}*/}
            {/*    onDismiss={() => setShowSnackbar(false)}*/}
            {/*    action={{*/}
            {/*        label: 'Got it',*/}
            {/*        onPress: () => setShowSnackbar(false),*/}
            {/*    }}*/}
            {/*>*/}
            {/*    No internet connection. Submitted data will be stored locally.*/}
            {/*</Snackbar>*/}
        </ScrollView>
    );
};
