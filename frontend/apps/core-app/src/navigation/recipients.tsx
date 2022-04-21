import * as React from 'react';
import { createStackNavigator } from '@react-navigation/stack';

import { RecipientListScreen } from '../screens/RecipientListScreen';
import { RecipientProfileScreen } from '../screens/RecipientProfileScreen';
import { RecipientRegistrationScreen } from '../screens/RecipientRegistrationScreen';
import { routes } from '../constants/routes';

const Stack = createStackNavigator();

export const RecipientNavigator: React.FC = () => {
  return (
    <Stack.Navigator initialRouteName={routes.recipientsList.name}>
      <Stack.Screen
        name={routes.recipientsList.name}
        component={RecipientListScreen}
        options={{
          title: routes.recipientsList.pageTitle,
          header: () => <></>,
        }}
      />
      <Stack.Screen
        name={routes.recipientsProfile.name}
        component={RecipientProfileScreen}
        options={{
          title: routes.recipientsProfile.pageTitle,
          headerTitle: routes.recipientsProfile.headerTitle,
        }}
      />
      <Stack.Screen
        name={routes.recipientsRegistration.name}
        component={RecipientRegistrationScreen}
        options={{
          title: routes.recipientsRegistration.pageTitle,
          headerTitle: routes.recipientsRegistration.headerTitle,
        }}
      />
    </Stack.Navigator>
  );
};
