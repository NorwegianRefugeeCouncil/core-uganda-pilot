import { Text } from 'react-native-paper';
import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { StyleProp, View } from 'react-native';

type BodyProps = {
  style: StyleProp<any>;
  children: React.ReactNode;
};

const Body: React.FC<BodyProps> = ({ style, children }) => {
  return (
    <View style={style}>
      {children}
    </View>
  );
}

export default Body;