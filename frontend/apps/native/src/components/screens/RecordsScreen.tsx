import React from 'react';
import {FAB, Title} from 'react-native-paper';
import {layout} from '../../styles';
import routes from '../../constants/routes';
import {FlatList, Text, TouchableOpacity, View} from 'react-native';
import useApiClient from "../../utils/clients";
import uuidv4 from 'uuid';
import {ScreenProps} from "../Router";
import {RECORD_ACTIONS} from "../../reducers/recordsReducers";

const RecordsScreen: React.FC<ScreenProps> = ({navigation, route, state, dispatch}) => {
    const {formId, databaseId} = route.params;
    const [isLoading, setIsLoading] = React.useState(true);

    const client = useApiClient();

    React.useEffect(() => {
        if (!formId){
            return
        }
        if (!databaseId){
            return
        }
        client.listRecords({formId, databaseId})
            .then((data) => {
                dispatch({
                    type: RECORD_ACTIONS.GET_RECORDS, payload: {
                        formId,
                        records: data.response?.items
                    }
                })
                setIsLoading(false)
            })
    }, [client, formId, databaseId]);

    return (
        <View style={[layout.container, layout.body]}>
            <Title>{routes.records.title}</Title>
            {!isLoading && (
                <View>
                    <FlatList
                        style={{width: '100%'}}
                        data={state.formsById[formId].records}
                        renderItem={({item}) => (
                            <TouchableOpacity
                                key={item.id}
                                onPress={() => navigation.navigate(routes.viewRecord.name, {recordId: item.id, formId})}
                            >
                                <View style={{flexDirection: 'row', flex: 1}}>
                                    <View style={{justifyContent: 'center', paddingRight: 12}}>
                                        <Text>{item.id}</Text>
                                    </View>
                                </View>
                            </TouchableOpacity>
                        )}
                    />
                    <FlatList
                        style={{width: '100%'}}
                        data={state.formsById[formId].localRecords}
                        renderItem={({item, index}) => (
                            <TouchableOpacity
                                key={index}
                                onPress={() => navigation.navigate(routes.addRecord.name, {recordId: item, formId})}
                            >
                                <View style={{flexDirection: 'row', flex: 1}}>
                                    <View style={{justifyContent: 'center', paddingRight: 12}}>
                                        <Text>{item}</Text>
                                    </View>
                                </View>
                            </TouchableOpacity>
                        )}
                    />
                </View>
            )}

            <FAB
                style={layout.fab}
                icon="plus"
                color={'white'}
                onPress={() => navigation.navigate(routes.addRecord.name, {formId, recordId: uuidv4()})}
            />
        </View>
    );
};

export default RecordsScreen;
