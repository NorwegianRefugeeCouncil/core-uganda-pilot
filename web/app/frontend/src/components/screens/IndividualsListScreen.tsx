import React from 'react';
import {FAB, Text, Title} from 'react-native-paper';
import {layout} from '../../styles';
import routes from '../../constants/routes';
import {FlatList, Image, TouchableOpacity, View} from 'react-native';
import pngIndividual from '../../../assets/png/symbol_individuals.png';
import theme from '../../constants/theme';

type Individual = {
    id: string,
    name: string
}

export const testIndividuals: Individual[] = [
    {
        id: 'c529d679-3bb6-4a20-8f06-c096f4d9adc1',
        name: 'Person 1'
    },
    {
        id: 'c529d679-3bb6-4a20-8f06-c096f4d9adc2',
        name: 'Person 2'
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
                        onPress={() => navigation.navigate(routes.individual.name, {id: index})}
                    >
                        <View style={{flexDirection: 'row', flex: 1}}>
                            <View style={{justifyContent: 'center', paddingRight: 12}}>
                                <Image source={pngIndividual}
                                       style={{tintColor: theme.colors.text, width: 20, height: 20}}/>
                            </View>
                            <View style={{justifyContent: 'center'}}>
                                <Title>{index}</Title>
                                <Text>{item.name}</Text>
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
