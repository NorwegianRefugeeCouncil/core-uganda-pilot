import React from 'react';
import { FAB, List, Title } from 'react-native-paper';
import { layout } from '../../styles';
import routes from '../../constants/routes';
import pngIndividuals from '../../../assets/png/symbol_individuals.png';
import theme from '../../constants/theme';
import { View } from 'react-native';

type Individual = {
  name: string
}
export const testIndividuals: Individual[] = [
  {
    name: 'Person 1'
  },
  {
    name: 'Person 2'
  }
];

const IndividualsListScreen: React.FC<any> = ({ navigation }) => {
  return (
    <View style={layout.body}>
      <Title>Individuals</Title>
      <FAB
        style={layout.fab}
        icon="plus"
        color={'white'}
        onPress={() => navigation.navigate(routes.individual.name, {id: null})}
      />

      <List.Section>
        {testIndividuals.map((individual, index) =>
          <List.Item
            key={index}
            title={individual.name}
            left={() => <List.Icon icon={pngIndividuals} color={theme.colors.backdrop}/>}
            onPress={() => navigation.navigate(routes.individual.name, {id: index})}
          />
        )}
      </List.Section>
    </View>
  );
};

export default IndividualsListScreen;