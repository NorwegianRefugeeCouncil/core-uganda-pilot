"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var styles_1 = require("../styles");
var react_1 = require("react");
var NavigationBar_1 = require("./NavigationBar");
var native_1 = require("@react-navigation/native");
var stack_1 = require("@react-navigation/stack");
var FormsScreen_1 = require("./screens/FormsScreen");
var routes_1 = require("../constants/routes");
var theme_1 = require("../constants/theme");
var host_1 = require("../constants/host");
var DesignSystemDemoScreen_1 = require("./screens/DesignSystemDemoScreen");
var RecordsScreen_1 = require("./screens/RecordsScreen");
var AddRecordScreen_1 = require("./screens/AddRecordScreen");
var ViewRecordScreen_1 = require("./screens/ViewRecordScreen");
var LoginCallbackScreen_1 = require("./screens/LoginCallbackScreen");
var recordsReducers_1 = require("../reducers/recordsReducers");
function Router() {
    var Stack = (0, stack_1.createStackNavigator)();
    var _a = (0, react_1.useReducer)(recordsReducers_1.recordsReducer, recordsReducers_1.initialRecordsState), state = _a[0], dispatch = _a[1];
    var linkingConfig = {
        prefixes: [host_1["default"]],
        config: {
            screens: {
                Forms: routes_1["default"].forms.name,
                Records: routes_1["default"].records.name,
                AddRecord: routes_1["default"].addRecord.name,
                ViewRecord: routes_1["default"].viewRecord.name,
                DesignSystem: routes_1["default"].designSystem.name,
                LoginCallback: 'callback'
            }
        }
    };
    return (<react_native_1.View style={styles_1.layout.container}>
            <native_1.NavigationContainer theme={theme_1.NavigationTheme} linking={linkingConfig}>
                <Stack.Navigator initialRouteName={routes_1["default"].forms.name}>
                    <Stack.Group screenOptions={{
            header: function (props) { return <NavigationBar_1["default"] {...props}/>; }
        }}>
                        <Stack.Screen name={routes_1["default"].forms.name} component={FormsScreen_1["default"]} options={{
            title: routes_1["default"].forms.title
        }}/>
                        <Stack.Screen name={'callback'} component={LoginCallbackScreen_1["default"]} options={{
            title: routes_1["default"].forms.title
        }}/>
                        <Stack.Screen name={routes_1["default"].records.name} options={{
            title: routes_1["default"].records.title
        }}>
                            {function (props) {
            return <RecordsScreen_1["default"] {...props} state={state} dispatch={dispatch}/>;
        }}
                        </Stack.Screen>
                        <Stack.Screen name={routes_1["default"].addRecord.name} options={{
            title: routes_1["default"].addRecord.title
        }}>
                            {function (props) {
            return <AddRecordScreen_1["default"] {...props} state={state} dispatch={dispatch}/>;
        }}
                        </Stack.Screen>
                        <Stack.Screen name={routes_1["default"].viewRecord.name} options={{
            title: routes_1["default"].viewRecord.title
        }}>
                            {function (props) {
            return <ViewRecordScreen_1["default"] {...props} state={state} dispatch={dispatch}/>;
        }}
                        </Stack.Screen>
                        <Stack.Screen name={routes_1["default"].designSystem.name} component={DesignSystemDemoScreen_1["default"]} options={{
            title: routes_1["default"].designSystem.title
        }}/>
                    </Stack.Group>
                </Stack.Navigator>
            </native_1.NavigationContainer>
        </react_native_1.View>);
}
exports["default"] = Router;
