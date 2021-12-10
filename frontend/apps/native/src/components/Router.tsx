import { View } from 'react-native';
import React, { useReducer } from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import Constants from 'expo-constants';

import { layout } from '../styles';
import routes from '../constants/routes';
import { NavigationTheme } from '../constants/theme';
import { initialRecordsState, recordsReducer, RecordsStoreProps } from '../reducers/recordsReducers';
import { AddRecordScreenContainerProps, ViewRecordScreenContainerProps } from '../types/screens';

import NavigationBar from './NavigationBar';
import FormsScreen from './screens/FormsScreen';
import DesignSystemDemoScreen from './screens/DesignSystemDemoScreen';
import RecordsScreen from './screens/RecordsScreen';
import ViewRecordScreen from './screens/ViewRecordScreen';
import { FormsScreenContainer } from './screen_containers/FormsScreenContainer';
import { AddRecordScreenContainer } from './screen_containers/AddRecordScreenContainer';
import { ViewRecordScreenContainer } from './screen_containers/ViewRecordScreenContainer';

export type ScreenProps = {
  navigation: any;
  route: any;
  state: RecordsStoreProps;
  dispatch: any;
};

export default function Router() {
  const Stack = createStackNavigator();
  const [state, dispatch] = useReducer(recordsReducer, initialRecordsState);

  const linkingConfig = {
    // TODO : revisit this
    prefixes: [Constants.manifest?.extra?.server_uri],
    config: {
      screens: {
        Forms: routes.forms.name,
        Records: routes.records.name,
        AddRecord: routes.addRecord.name,
        ViewRecord: routes.viewRecord.name,
        DesignSystem: routes.designSystem.name,
        LoginCallback: 'callback',
      },
    },
  };

  return (
    <View style={layout.container}>
      <NavigationContainer theme={NavigationTheme} linking={linkingConfig}>
        <Stack.Navigator initialRouteName={routes.forms.name}>
          <Stack.Group
            screenOptions={{
              header: (props) => <NavigationBar {...props} />,
            }}
          >
            <Stack.Screen
              name={routes.forms.name}
              component={FormsScreenContainer}
              options={{
                title: routes.forms.title,
              }}
            />
            <Stack.Screen
              name={routes.records.name}
              options={{
                title: routes.records.title,
              }}
            >
              {(props) => <RecordsScreen {...props} state={state} dispatch={dispatch} />}
            </Stack.Screen>
            <Stack.Screen
              name={routes.addRecord.name}
              options={{
                title: routes.addRecord.title,
              }}
            >
              {({ navigation, route }) => (
                <AddRecordScreenContainer
                  navigation={navigation}
                  route={route as AddRecordScreenContainerProps['route']}
                  state={state}
                  dispatch={dispatch}
                />
              )}
            </Stack.Screen>
            <Stack.Screen
              name={routes.viewRecord.name}
              options={{
                title: routes.viewRecord.title,
              }}
            >
              {({ navigation, route }) => (
                <ViewRecordScreenContainer
                  navigation={navigation}
                  route={route as ViewRecordScreenContainerProps['route']}
                  state={state}
                  dispatch={dispatch}
                />
              )}
            </Stack.Screen>
            <Stack.Screen
              name={routes.designSystem.name}
              component={DesignSystemDemoScreen}
              options={{
                title: routes.designSystem.title,
              }}
            />
          </Stack.Group>
        </Stack.Navigator>
      </NavigationContainer>
    </View>
  );
}
