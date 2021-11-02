import {View} from 'react-native';
import {layout} from '../styles';
import React from 'react';
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

export default function Router() {
    const Stack = createStackNavigator();

    const linkingConfig = {
        prefixes: [host],
        config: {
            screens: {
                Forms: routes.forms.name,
                Records: routes.records.name,
                Record: routes.addRecord.name,
                DesignSystem: routes.designSystem.name
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
                                title: routes.forms.title
                            }}
                        />
                        <Stack.Screen
                            name={routes.records.name}
                            component={RecordsScreen}
                            options={{
                                title: routes.records.title
                            }}
                        />
                        <Stack.Screen
                            name={routes.addRecord.name}
                            component={AddRecordScreen}
                            options={{
                                title: routes.addRecord.title
                            }}
                        />
                        <Stack.Screen
                            name={routes.viewRecord.name}
                            component={ViewRecordScreen}
                            options={{
                                title: routes.viewRecord.title
                            }}
                        />
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

