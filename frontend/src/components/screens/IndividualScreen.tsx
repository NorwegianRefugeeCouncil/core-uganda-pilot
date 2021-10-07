import React from 'react';
import { FAB, Text } from 'react-native-paper';
import { common, layout } from '../../styles';
import { testIndividuals } from './IndividualsListScreen';
import { darkTheme } from '../../constants/theme';
import { View } from 'react-native';

const IndividualScreen: React.FC<any> = ({ route }) => {
  const { id } = route.params;

  return (
    <View style={[layout.container, layout.body, common.darkBackground]}>
      <Text theme={darkTheme}>
        {id == null ? 'new person' : testIndividuals[id].name}
      </Text>
      <FAB
        style={layout.fab}
        icon="chevron-right"
        color={darkTheme.colors.white}
        onPress={() => console.log('Pressed')}
      />
    </View>
  );
};

export default IndividualScreen;