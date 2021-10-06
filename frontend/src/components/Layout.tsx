import { View } from 'react-native';
import layout from '../styles/layout';
import React from 'react';
import Navigation from './Navigation';
import Body from './Body';

export default function Layout() {
  return (
    <View style={layout.container}>
      <Navigation style={layout.navigation}/>
      <Body style={layout.body} />
    </View>
  );
}