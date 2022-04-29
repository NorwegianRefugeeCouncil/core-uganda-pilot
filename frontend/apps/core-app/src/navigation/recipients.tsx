import * as React from 'react';
import { createStackNavigator } from '@react-navigation/stack';

import { RecipientListScreen } from '../screens/RecipientListScreen';
import { RecipientProfileScreen } from '../screens/RecipientProfileScreen';
import { RecipientRegistrationScreen } from '../screens/RecipientRegistrationScreen';
import { routes } from '../constants/routes';

export type RecipientNavigatorParamList = {
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

const Stack = createStackNavigator<RecipientNavigatorParamList>();

export const RecipientNavigator: React.FC = () => {
  return (
    <Stack.Navigator initialRouteName={routes.recipientsList.name}>
      <Stack.Screen
        name={routes.recipientsList.name}
        component={RecipientListScreen}
        options={{ headerShown: false }}
      />
      <Stack.Screen
        name={routes.recipientsProfile.name}
        component={RecipientProfileScreen}
        options={{ headerShown: false }}
      />
      <Stack.Screen
        name={routes.recipientsRegistration.name}
        component={RecipientRegistrationScreen}
        options={{ headerShown: false }}
      />
    </Stack.Navigator>
  );
};
