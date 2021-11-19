import {View} from 'react-native';
import {layout} from '../styles';
import React, {useReducer} from 'react';
import NavigationBar from './NavigationBar';
import {NavigationContainer} from '@react-navigation/native';
import {createStackNavigator} from '@react-navigation/stack';
import FormsScreen from './screens/FormsScreen';
import routes from '../constants/routes';
import {NavigationTheme} from '../constants/theme';
import host from "../constants/host";
import DesignSystemDemoScreen from "./screens/DesignSystemDemoScreen";
import RecordsScreen from "./screens/RecordsScreen";
import AddRecordScreen from "./screens/AddRecordScreen";
import ViewRecordScreen from "./screens/ViewRecordScreen";
import LoginCallbackScreen from "./screens/LoginCallbackScreen";
import {initialRecordsState, recordsReducer, RecordsStoreProps} from "../reducers/recordsReducers";

export type ScreenProps = {
    navigation: any,
    route: any,
    state: RecordsStoreProps,
    dispatch: any
}

export default function Router() {
    const Stack = createStackNavigator();
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
                LoginCallback: 'callback'
            }
        }
    };

    return (
        <View style={layout.container}>
            <NavigationContainer theme={NavigationTheme} linking={linkingConfig}>
                <Stack.Navigator initialRouteName={routes.forms.name}>
                    <Stack.Group
                        screenOptions={{
                            header: (props) => <NavigationBar {...props} />
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
                            name={'callback'}
                            component={LoginCallbackScreen}
                            options={{
                                title: routes.forms.title,
                            }}
                        />
                        <Stack.Screen
                            name={routes.records.name}
                            options={{
                                title: routes.records.title
                            }}
                        >
                            {(props) =>
                                <RecordsScreen {...props} state={state} dispatch={dispatch}/>
                            }
                        </Stack.Screen>
                        <Stack.Screen
                            name={routes.addRecord.name}
                            options={{
                                title: routes.addRecord.title
                            }}
                        >
                            {(props) =>
                                <AddRecordScreen {...props} state={state} dispatch={dispatch}/>
                            }
                        </Stack.Screen>
                        <Stack.Screen
                            name={routes.viewRecord.name}
                            options={{
                                title: routes.viewRecord.title
                            }}
                        >
                            {(props) =>
                                <ViewRecordScreen {...props} state={state} dispatch={dispatch}/>
                            }
                        </Stack.Screen>
                        <Stack.Screen
                            name={routes.designSystem.name}
                            component={DesignSystemDemoScreen}
                            options={{
                                title: routes.designSystem.title
                            }}
                        />
                    </Stack.Group>
                </Stack.Navigator>
            </NavigationContainer>
        </View>
    );
}
