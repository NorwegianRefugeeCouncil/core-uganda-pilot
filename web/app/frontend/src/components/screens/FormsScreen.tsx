import React from 'react';
import {Title} from 'react-native-paper';
import {layout} from '../../styles';
import {FlatList, Text, TouchableOpacity, View} from 'react-native';
import useApiClient from "../../utils/clients";
import routes from "../../constants/routes";
import {FormDefinition} from "core-js-api-client/lib/types/types";

const FormsScreen: React.FC<any> = ({navigation}) => {
    const [forms, setForms] = React.useState<FormDefinition[]>();
    const [isLoading, setIsLoading] = React.useState(true);
    const client = useApiClient();

    React.useEffect(() => {
        client.listForms({})
            .then((data) => {
                setForms(data.response?.items)
                setIsLoading(false)
            })
    }, [client]);

    return (
        <View style={layout.body}>
            <Title>{routes.forms.title}</Title>
            {!isLoading && (
                <FlatList
                    style={{flex: 1, width: '100%'}}
                    data={forms}
                    renderItem={({item, index, separators}) => (
                        <TouchableOpacity
                            key={index}
                            onPress={() => navigation.navigate(routes.records.name, {
                                formId: item.id,
                                databaseId: item.databaseId
                            })}
                        >
                            <View style={{flexDirection: 'row', flex: 1}}>
                                <View style={{justifyContent: 'center', paddingRight: 12}}>
                                    <Text>{item.code}</Text>
                                </View>
                                <View style={{justifyContent: 'center'}}>
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

export default FormsScreen;
