import * as React from 'react';
import { createStackNavigator } from '@react-navigation/stack';

import { RecipientListScreen } from '../screens/RecipientListScreen';
import { RecipientProfileScreen } from '../screens/RecipientProfileScreen';

const Stack = createStackNavigator();

export const RecipientNavigator: React.FC = () => {
  return (
    <Stack.Navigator initialRouteName="RecipientList">
      <Stack.Screen
        name="RecipientList"
        component={RecipientListScreen}
        options={{
          title: 'Recipient List',
          header: () => <></>,
        }}
      />
      <Stack.Screen
        name="RecipientProfile"
        component={RecipientProfileScreen}
        options={{
          title: 'Recipient Profile',
          headerTitle: 'Recipient Profile',
        }}
      />
    </Stack.Navigator>
  );
};
