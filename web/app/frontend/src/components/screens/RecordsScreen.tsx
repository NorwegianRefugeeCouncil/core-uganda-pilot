import React from 'react';
import {FAB, Title} from 'react-native-paper';
import {layout} from '../../styles';
import routes from '../../constants/routes';
import {FlatList, Text, TouchableOpacity, View} from 'react-native';
import useApiClient from "../../utils/clients";
import uuidv4 from 'uuid';

const RecordsScreen: React.FC<any> = ({navigation, route}) => {
    const [records, setRecords] = React.useState<any>();
    const [isLoading, setIsLoading] = React.useState(true);
    const client = useApiClient();
    const {id, databaseId} = route.params;

    React.useEffect(() => {
        client.listRecords({formId: id, databaseId})
            .then((data) => {
                setRecords(data.response?.items)
                setIsLoading(false)
            })
    }, []);
    // console.log('RECORDS', records)

    return (
        <View style={layout.body}>
            <Title>{routes.records.title}</Title>
            {!isLoading && (
                <FlatList
                    style={{flex: 1, width: '100%'}}
                    data={records}
                    renderItem={({item, index, separators}) => (
                        <TouchableOpacity
                            key={index}
                            onPress={() => navigation.navigate(routes.viewRecord.name, {id: item.id})}
                        >
                            <View style={{flexDirection: 'row', flex: 1}}>
                                <View style={{justifyContent: 'center', paddingRight: 12}}>
                                    <Text>{item.id}</Text>
                                </View>
                            </View>
                        </TouchableOpacity>
                    )}
                />
            )}
            <FAB
                style={layout.fab}
                icon="plus"
                color={'white'}
                onPress={() => navigation.navigate(routes.addRecord.name, {formId: id, recordId: uuidv4()})}
            />
        </View>
    );
};

export default RecordsScreen;
