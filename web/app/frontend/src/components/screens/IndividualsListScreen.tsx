import React from 'react';
import {FAB, Title} from 'react-native-paper';
import {layout} from '../../styles';
import routes from '../../constants/routes';
import {FlatList, Image, TouchableOpacity, View, Text} from 'react-native';
import pngIndividual from '../../../assets/png/symbol_individuals.png';
import theme from '../../constants/theme';
import iamClient from "../../utils/clients";
import {PartyAttributeDefinitionList} from "core-js-api-client/lib/types/models";
import {Individual} from "../../../../client/src/types/models";
import _ from "lodash";
import {Subject} from "rxjs";
import {FlatIndividual} from "./IndividualScreen";

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
    const [individuals, setIndividuals] = React.useState<Individual[]>();
    let individualsSubject = new Subject([]);

    React.useEffect(() => {
        individualsSubject.pipe(iamClient.Individuals().List()).subscribe(
            (data: { items: Individual[] }) => {
                setIndividuals(data.items)
            }
        );
        individualsSubject.next();

        return () => {
            if (individualsSubject) {
                individualsSubject.unsubscribe();
            }
        };
    }, []);

    return (
        <View style={layout.body}>
            <Title>Individuals</Title>
            <FlatList
                style={{flex: 1, width: '100%'}}
                data={individuals}
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
                                <Text>{item.attributes['8514da51-aad5-4fb4-a797-8bcc0c969b27']}</Text>
                            </View>
                        </View>
                    </TouchableOpacity>
                )}
            />
            <FAB
                style={layout.fab}
                icon="plus"
                color={'white'}
                onPress={() => navigation.navigate(routes.individual.name, {id: 'new'})}
            />
        </View>
    );
};

export default IndividualsListScreen;
