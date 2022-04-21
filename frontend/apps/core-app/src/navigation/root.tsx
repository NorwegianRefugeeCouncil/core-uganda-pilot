import * as React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';

import { routes } from '../constants/routes';

import { RecipientNavigator } from './recipients';

const Drawer = createDrawerNavigator();

export const RootNavigator: React.FC = () => {
  return (
    <Drawer.Navigator initialRouteName={routes.recipientsRoot.name}>
      <Drawer.Screen
        name={routes.recipientsRoot.name}
        component={RecipientNavigator}
        options={{
          title: routes.recipientsRoot.pageTitle,
          headerTitle: routes.recipientsRoot.headerTitle,
        }}
      />
    </Drawer.Navigator>
  );
};
