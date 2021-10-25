import React from 'react';
import {FAB, Title} from 'react-native-paper';
import {layout} from '../../styles';
import routes from '../../constants/routes';
import {FlatList, Image, TouchableOpacity, View, Text} from 'react-native';
import pngIndividual from '../../../assets/png/symbol_individuals.png';
import theme from '../../constants/theme';

type TestIndividual = {
    id: string
}

export const testIndividuals: TestIndividual[] = [
    {
        id: 'c529d679-3bb6-4a20-8f06-c096f4d9adc1'
    },
    {
        id: 'bbf539fd-ebaa-4438-ae4f-8aca8b327f42'
    }
];

const IndividualsListScreen: React.FC<any> = ({navigation}) => {
    return (
        <View style={layout.body}>
            <Title>Individuals</Title>
            <FlatList
                style={{flex: 1, width: '100%'}}
                data={testIndividuals}
                renderItem={({item, index, separators}) => (
                    <TouchableOpacity
                        key={index}
                        onPress={() => navigation.navigate(routes.individual.name, {id: item.id})}
                    >
                        <View style={{flexDirection: 'row', flex: 1}}>
                            <View style={{justifyContent: 'center', paddingRight: 12}}>
                                <Image source={pngIndividual}
                                       style={{tintColor: theme.colors.text, width: 20, height: 20}}/>
                            </View>
                            <View style={{justifyContent: 'center'}}>
                                <Title>{index}</Title>
                                <Text>{item.id}</Text>
                            </View>
                        </View>
                    </TouchableOpacity>
                )}
            />
            <FAB
                style={layout.fab}
                icon="plus"
                color={'white'}
                onPress={() => navigation.navigate(routes.individual.name, {id: null})}
            />
        </View>
    );
};

export default IndividualsListScreen;
