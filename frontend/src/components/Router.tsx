import { View } from 'react-native';
import { layout } from '../styles';
import React from 'react';
import NavigationBar from './NavigationBar';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import IndividualsListScreen from './screens/IndividualsListScreen';
import HomeScreen from './screens/HomeScreen';
import IndividualScreen from './screens/IndividualScreen';
import routes from '../constants/routes';
import CasesScreen from './screens/CasesScreen';
import { NavigationTheme } from '../constants/theme';


export default function Router() {
  const Stack = createStackNavigator();

  return (
    <View style={layout.container}>
      <NavigationContainer theme={NavigationTheme}>
        <Stack.Navigator initialRouteName={routes.home.name}>
          <Stack.Group
            screenOptions={{
              header: (props) => <NavigationBar {...props} />
            }}
          >
            <Stack.Screen
              name={routes.home.name}
              component={HomeScreen}
            />
            <Stack.Screen
              name={routes.individuals.name}
              component={IndividualsListScreen}
            />
            <Stack.Screen
              name={routes.cases.name}
              component={CasesScreen}
            />
          </Stack.Group>
          <Stack.Group screenOptions={{ presentation: 'modal' }}>
            <Stack.Screen
              name={routes.individual.name}
              component={IndividualScreen}
            />
          </Stack.Group>
        </Stack.Navigator>
      </NavigationContainer>
    </View>
  );
}

