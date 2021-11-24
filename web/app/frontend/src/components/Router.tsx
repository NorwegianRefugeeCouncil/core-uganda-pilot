import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import React, { useReducer } from 'react';
import { View } from 'react-native';

import host from '../constants/host';
import routes from '../constants/routes';
import { NavigationTheme } from '../constants/theme';
import {
    initialRecordsState,
    recordsReducer,
} from '../reducers/recordsReducers';
import { layout } from '../styles';
import {
    AddRecordScreenProps,
    RecordsScreenProps,
    StackParamList,
    ViewRecordScreenProps,
} from '../types/screens';
import NavigationBar from './NavigationBar';
import AddRecordScreen from './screens/AddRecordScreen';
import DesignSystemDemoScreen from './screens/DesignSystemDemoScreen';
import FormsScreen from './screens/FormsScreen';
import RecordsScreen from './screens/RecordsScreen';
import ViewRecordScreen from './screens/ViewRecordScreen';

export default function Router() {
    const Stack = createStackNavigator<StackParamList>();
    const [state, dispatch] = useReducer(recordsReducer, initialRecordsState);

    const linkingConfig = {
        prefixes: [host],
        config: {
            screens: {
                Forms: routes.forms.name,
                Records: routes.records.name,
                AddRecord: routes.addRecord.name,
                ViewRecord: routes.viewRecord.name,
                DesignSystem: routes.designSystem.name,
            },
        },
    };

    return (
        <View style={layout.container}>
            <NavigationContainer
                theme={NavigationTheme}
                linking={linkingConfig}
            >
                <Stack.Navigator initialRouteName={routes.forms.name}>
                    <Stack.Group
                        screenOptions={{
                            header: props => <NavigationBar {...props} />,
                        }}
                    >
                        <Stack.Screen
                            name={routes.forms.name}
                            component={FormsScreen}
                            options={{
                                title: routes.forms.title,
                            }}
                        />
                        <Stack.Screen
                            name={routes.records.name}
                            initialParams={{ formId: '', databaseId: '' }}
                            options={{
                                title: routes.records.title,
                            }}
                        >
                            {({ navigation, route }) => (
                                <RecordsScreen
                                    navigation={navigation}
                                    route={route as RecordsScreenProps['route']}
                                    state={state}
                                    dispatch={dispatch}
                                />
                            )}
                        </Stack.Screen>
                        <Stack.Screen
                            name={routes.addRecord.name}
                            options={{
                                title: routes.addRecord.title,
                            }}
                        >
                            {({ navigation, route }) => (
                                <AddRecordScreen
                                    navigation={navigation}
                                    route={
                                        route as AddRecordScreenProps['route']
                                    }
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
                                <ViewRecordScreen
                                    navigation={navigation}
                                    route={
                                        route as ViewRecordScreenProps['route']
                                    }
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
