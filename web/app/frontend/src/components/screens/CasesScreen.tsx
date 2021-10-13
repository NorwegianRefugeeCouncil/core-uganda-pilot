import React from 'react';
import { FAB, Text, Title } from 'react-native-paper';
import { layout } from '../../styles';
import routes from '../../constants/routes';
import { View } from 'react-native';

const CasesScreen: React.FC<any> = ({ navigation }) => {
  return (
    <View style={layout.body}>
        <Title>Cases</Title>
    </View>
  );
};

export default CasesScreen;