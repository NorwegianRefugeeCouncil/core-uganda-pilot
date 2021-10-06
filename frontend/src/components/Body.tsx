import { Text } from 'react-native-paper';
import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { StyleProp, View } from 'react-native';

type BodyProps = {
  style: StyleProp<any>;
};

const Body: React.FC<BodyProps> = ({ style }) => {
  return (
    <View style={style}>
      <Text>Open up App.tsx to start working on your app!</Text>
      <StatusBar style="auto" />
    </View>
  );
}

export default Body;