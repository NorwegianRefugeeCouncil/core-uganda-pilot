import * as React from 'react';
import { View, Text } from 'react-native';
import { createDrawerNavigator } from '@react-navigation/drawer';

const Drawer = createDrawerNavigator();

export const RootNavigator: React.FC = () => {
  return (
    <Drawer.Navigator>
      <Drawer.Screen
        name="Screen one"
        component={() => (
          <View>
            <Text>Page one</Text>
          </View>
        )}
      />
      <Drawer.Screen
        name="Screen two"
        component={() => (
          <View>
            <Text>Page two</Text>
          </View>
        )}
      />
      <Drawer.Screen
        name="Screen three"
        component={() => (
          <View>
            <Text>Page three</Text>
          </View>
        )}
      />
    </Drawer.Navigator>
  );
};
