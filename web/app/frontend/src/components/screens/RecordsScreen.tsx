import React from 'react';
import { FlatList, Text, TouchableOpacity, View } from 'react-native';

import routes from '../../constants/routes';
import {RECORD_ACTIONS, RecordsStoreProps} from '../../reducers/recordsReducers';
import { layout } from '../../styles';
import {RecordsScreenContainerProps, StackParamList} from '../../types/screens';
import { useApiClient } from '../../utils/useApiClient';
import {StackNavigationProp} from "@react-navigation/stack";
import {RecordInstance} from "core-js-api-client";
import {RouteProp} from "@react-navigation/native";

export type RecordsScreenProps = {
    navigation: StackNavigationProp<StackParamList, 'records'>;
    isLoading: boolean;
    records: RecordInstance[];
    localRecords: string[];
    formId: string;

}

export const RecordsScreen = ({
    navigation,
    isLoading,
    records,
    localRecords,
    formId,
}: RecordsScreenProps) => {

    return (
        <View style={[layout.container, layout.body]}>
            {/*<Title>{routes.records.title}</Title>*/}
            {!isLoading && (
                <View>
                    <FlatList
                        style={{ width: '100%' }}
                        data={records}
                        renderItem={({ item }) => (
                            <TouchableOpacity
                                key={item.id}
                                onPress={() =>
                                    navigation.navigate(
                                        routes.viewRecord.name,
                                        { recordId: item.id, formId }
                                    )
                                }
                            >
                                <View style={{ flexDirection: 'row', flex: 1 }}>
                                    <View
                                        style={{
                                            justifyContent: 'center',
                                            paddingRight: 12,
                                        }}
                                    >
                                        <Text>{item.id}</Text>
                                    </View>
                                </View>
                            </TouchableOpacity>
                        )}
                    />
                    <FlatList
                        style={{ width: '100%' }}
                        data={localRecords}
                        renderItem={({ item, index }) => (
                            <TouchableOpacity
                                key={index}
                                onPress={() =>
                                    navigation.navigate(routes.addRecord.name, {
                                        recordId: item,
                                        formId,
                                    })
                                }
                            >
                                <View style={{ flexDirection: 'row', flex: 1 }}>
                                    <View
                                        style={{
                                            justifyContent: 'center',
                                            paddingRight: 12,
                                        }}
                                    >
                                        <Text>{item}</Text>
                                    </View>
                                </View>
                            </TouchableOpacity>
                        )}
                    />
                </View>
            )}

            {/*<FAB*/}
            {/*    style={layout.fab}*/}
            {/*    icon="plus"*/}
            {/*    color={'white'}*/}
            {/*    onPress={() =>*/}
            {/*        navigation.navigate(routes.addRecord.name, {*/}
            {/*            formId,*/}
            {/*            recordId: uuidv4(),*/}
            {/*        })*/}
            {/*    }*/}
            {/*/>*/}
        </View>
    );
};
