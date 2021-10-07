import React from 'react';
import { Text } from 'react-native-paper';
import layout from '../../styles/layout';
import Body from '../Body';

const HomeScreen = () => {
  return (
    <Body style={layout.body}>
      <Text>Home</Text>
    </Body>
  );
};

export default HomeScreen;