import * as React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';
import { createStackNavigator } from '@react-navigation/stack';
import { useBreakpointValue } from 'native-base';

import { routes } from '../constants/routes';
import { RecipientListScreen } from '../screens/RecipientListScreen';
import { RecipientProfileScreen } from '../screens/RecipientProfileScreen';
import { RecipientRegistrationScreen } from '../screens/RecipientRegistrationScreen';
import { LargeNavHeader } from '../components/NavHeader';

export type RootNavigatorParamList = {
  recipientsList: undefined;
  recipientsRegistration: {
    formId: string;
    databaseId: string;
  };
  recipientsProfile: {
    recordId: string;
    formId: string;
    databaseId: string;
  };
};

const Drawer = createDrawerNavigator<RootNavigatorParamList>();
const Stack = createStackNavigator<RootNavigatorParamList>();

const makeScreens = (Screen: typeof Drawer.Screen | typeof Stack.Screen) => (
  <>
    <Screen name={routes.recipientsList.name} component={RecipientListScreen} />
    <Screen
      name={routes.recipientsProfile.name}
      component={RecipientProfileScreen}
    />
    <Screen
      name={routes.recipientsRegistration.name}
      component={RecipientRegistrationScreen}
    />
  </>
);

const SmallRootNavigator: React.FC = () => {
  return (
    <Drawer.Navigator initialRouteName={routes.recipientsList.name}>
      {makeScreens(Drawer.Screen)}
    </Drawer.Navigator>
  );
};

const LargeRootNavigator: React.FC = () => {
  return (
    <Stack.Navigator
      initialRouteName={routes.recipientsList.name}
      screenOptions={{ header: LargeNavHeader }}
    >
      {makeScreens(Stack.Screen)}
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
