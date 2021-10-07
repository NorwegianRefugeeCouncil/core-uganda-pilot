import { View } from 'react-native';
import layout from '../styles/layout';
import React from 'react';
import Navigation from './Navigation';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import IndividualsScreen from './screens/IndividualsScreen';
import HomeScreen from './screens/HomeScreen';

const Stack = createStackNavigator();

export default function Layout() {
  return (
    <View style={layout.container}>
      <NavigationContainer>
        <Stack.Navigator
          initialRouteName="home"
          screenOptions={{
            header: (props) => <Navigation {...props}/>
          }}>
            <Stack.Screen name="home" component={HomeScreen}/>
            <Stack.Screen name="individuals" component={IndividualsScreen}/>
        </Stack.Navigator>
      </NavigationContainer>
    </View>
  );
}

