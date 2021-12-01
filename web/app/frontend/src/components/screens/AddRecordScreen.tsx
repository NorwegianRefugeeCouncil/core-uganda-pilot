import React from "react";
import { common, layout } from "../../styles";
import { Button, ScrollView, Text, View } from "react-native";
import FormControl from "../form/FormControl";
import { FormDefinition } from "core-js-api-client";
import { Control, FieldValues, FormState } from "react-hook-form";

export type AddRecordScreenProps = {
    form?: FormDefinition;
    control: Control<FieldValues, object>;
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
    isWeb, // TODO remove this?
    hasLocalData,
    isConnected,
    isLoading,
}: AddRecordScreenProps)=> {
  // console.log(form?.fields)
    return (
        <ScrollView contentContainerStyle={[layout.container, layout.body, common.darkBackground]}>
            <View style={[]}>
                {/* simulate network changes, for testing */}
                {!isWeb && (
                    <View style={{display: "flex", flexDirection: "row"}}>
                        {/*<Switch
                            value={simulateOffline}
                            onValueChange={() => {
                                setSimulateOffline(!simulateOffline)
                                setIsConnected(simulateOffline)
                                setShowSnackbar(!simulateOffline)
                            }}
                        />*/}
                        <Text> simulate being offline </Text>
                    </View>
                )}

                {/* upload data collected offline */}
                {hasLocalData && (
                    <View style={{ display: "flex", flexDirection: "column" }}>
                        <Text>There is locally stored data for this individual.</Text>
                    </View>
                )}
                {hasLocalData && isConnected && (
                    <View style={{display: "flex", flexDirection: "column"}}>
                        <Text>
                            Do you want to upload it?
                        </Text>
                        <Button
                            accessibilityLabel="Submit local data"
                            title="Submit local data"
                            onPress={onSubmit}
                        />
                    </View>
                )}
                {isLoading ? (
                    <Text>Loading...</Text>
                ) : (
                    <View style={{ width: "100%" }}>
                        {form?.fields.map(field => {
                            return (
                                <FormControl
                                    key={field.code}
                                    fieldDefinition={field}
                                    style={{ width: "100%" }}
                                    // value={''} // take value from record
                                    control={control}
                                    name={field.id}
                                    errors={formState?.errors}
                                />
                            );
                        })}
                        <Button accessibilityLabel="Submit" title="Submit" onPress={onSubmit} />
                    </View>
                )}
            </View>
        </ScrollView>
    );
};

export default AddRecordScreen;
