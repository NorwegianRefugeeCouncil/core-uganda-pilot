import React from 'react';
import { Title } from 'react-native-paper';
import { layout } from '../../styles';
import { View } from 'react-native';

const HomeScreen = () => {
  return (
    <View style={layout.body}>
      <Title>Home</Title>
    </View>
  );
};

export default HomeScreen;