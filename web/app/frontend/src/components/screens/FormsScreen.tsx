import { FormDefinition } from "core-js-api-client/lib/types/types";
import React from "react";
import { FlatList, Text, TouchableOpacity, View } from "react-native";

import routes from "../../constants/routes";
import { layout } from "../../styles";
import testIds from "../../constants/testIds";
import { StackParamList } from "../../types/screens";
import { StackNavigationProp } from "@react-navigation/stack";

export type FormsScreenProps = {
    isLoading: boolean;
    forms?: FormDefinition[];
    navigation: StackNavigationProp<StackParamList, "forms">;
};

export const FormsScreen = ({
    isLoading,
    forms,
    navigation,
}: FormsScreenProps) => {
    return (
        <View style={layout.body}>
            {/*<Title>{routes.forms.title}</Title>*/}
            {!isLoading && (
                <FlatList
                    style={{ flex: 1, width: "100%" }}
                    data={forms}
                    renderItem={({ item, index, separators }) => (
                        <TouchableOpacity
                            key={index}
                            onPress={() =>
                                navigation.navigate(
                                    routes.records.name as keyof StackParamList,
                                    {
                                        formId: item.id,
                                        databaseId: item.databaseId,
                                    }
                                )
                            }
                        >
                            <View
                                testID={testIds.formListItem}
                                style={{ flexDirection: "row", flex: 1 }}
                            >
                                <View
                                    style={{
                                        justifyContent: "center",
                                        paddingRight: 12,
                                    }}
                                >
                                    <Text>{item.code}</Text>
                                </View>
                                <View style={{ justifyContent: "center" }}>
                                    <Text>{item.name}</Text>
                                </View>
                            </View>
                        </TouchableOpacity>
                    )}
                />
            )}
        </View>
    );
};
