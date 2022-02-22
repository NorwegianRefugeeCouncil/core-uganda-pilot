import * as React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';

import { RecipientRegistrationScreen } from '../screens/RecipientRegistrationScreen';

import { RecipientNavigator } from './recipients';

const Drawer = createDrawerNavigator();

export const RootNavigator: React.FC = () => {
  return (
    <Drawer.Navigator initialRouteName="Recipients">
      <Drawer.Screen
        name="Recipients"
        component={RecipientNavigator}
        options={{
          title: 'Recipient List',
          headerTitle: 'Recipient List',
        }}
      />
      <Drawer.Screen
        name="RecipientRegistration"
        component={RecipientRegistrationScreen}
        options={{
          title: 'Recipient Registration',
          headerTitle: 'Recipient Registration',
        }}
      />
    </Drawer.Navigator>
  );
};
