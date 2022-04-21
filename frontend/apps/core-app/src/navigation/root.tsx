import * as React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';
import { createStackNavigator } from '@react-navigation/stack';
import { useBreakpointValue } from 'native-base';

import { routes } from '../constants/routes';
import { LargeNavHeader } from '../components/NavHeader';

import { RecipientNavigator } from './recipients';

const Drawer = createDrawerNavigator();

const SmallRootNavigator: React.FC = () => {
  return (
    <Drawer.Navigator initialRouteName={routes.recipientsRoot.name}>
      <Drawer.Screen
        name={routes.recipientsRoot.name}
        component={RecipientNavigator}
        options={{ title: routes.recipientsRoot.title }}
      />
    </Drawer.Navigator>
  );
};

const Stack = createStackNavigator();

const LargeRootNavigator: React.FC = () => {
  return (
    <Stack.Navigator
      initialRouteName={routes.recipientsRoot.name}
      screenOptions={{ header: LargeNavHeader }}
    >
      <Stack.Screen
        name={routes.recipientsRoot.name}
        component={RecipientNavigator}
        options={{ title: routes.recipientsRoot.title }}
      />
    </Stack.Navigator>
  );
};

export const RootNavigator: React.FC = () => {
  const Navigator = useBreakpointValue({
    base: SmallRootNavigator,
    sm: LargeRootNavigator,
  });

  return <Navigator />;
};
