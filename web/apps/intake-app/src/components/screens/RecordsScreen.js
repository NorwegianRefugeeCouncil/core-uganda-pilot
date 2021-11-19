"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var styles_1 = require("../../styles");
var routes_1 = require("../../constants/routes");
var react_native_1 = require("react-native");
var clients_1 = require("../../utils/clients");
var uuid_1 = require("uuid");
var recordsReducers_1 = require("../../reducers/recordsReducers");
var RecordsScreen = function (_a) {
    var navigation = _a.navigation, route = _a.route, state = _a.state, dispatch = _a.dispatch;
    var _b = route.params, formId = _b.formId, databaseId = _b.databaseId;
    var _c = react_1["default"].useState(true), isLoading = _c[0], setIsLoading = _c[1];
    var client = (0, clients_1["default"])();
    react_1["default"].useEffect(function () {
        client.listRecords({ formId: formId, databaseId: databaseId })
            .then(function (data) {
            var _a;
            dispatch({
                type: recordsReducers_1.RECORD_ACTIONS.GET_RECORDS, payload: {
                    formId: formId,
                    records: (_a = data.response) === null || _a === void 0 ? void 0 : _a.items
                }
            });
            setIsLoading(false);
        });
    }, [client]);
    return (<react_native_1.View style={[styles_1.layout.container, styles_1.layout.body]}>
            <react_native_paper_1.Title>{routes_1["default"].records.title}</react_native_paper_1.Title>
            {!isLoading && (<react_native_1.View>
                    <react_native_1.FlatList style={{ width: '100%' }} data={state.formsById[formId].records} renderItem={function (_a) {
                var item = _a.item;
                return (<react_native_1.TouchableOpacity key={item.id} onPress={function () { return navigation.navigate(routes_1["default"].viewRecord.name, { recordId: item.id, formId: formId }); }}>
                                <react_native_1.View style={{ flexDirection: 'row', flex: 1 }}>
                                    <react_native_1.View style={{ justifyContent: 'center', paddingRight: 12 }}>
                                        <react_native_1.Text>{item.id}</react_native_1.Text>
                                    </react_native_1.View>
                                </react_native_1.View>
                            </react_native_1.TouchableOpacity>);
            }}/>
                    <react_native_1.FlatList style={{ width: '100%' }} data={state.formsById[formId].localRecords} renderItem={function (_a) {
                var item = _a.item, index = _a.index;
                return (<react_native_1.TouchableOpacity key={index} onPress={function () { return navigation.navigate(routes_1["default"].addRecord.name, { recordId: item, formId: formId }); }}>
                                <react_native_1.View style={{ flexDirection: 'row', flex: 1 }}>
                                    <react_native_1.View style={{ justifyContent: 'center', paddingRight: 12 }}>
                                        <react_native_1.Text>{item}</react_native_1.Text>
                                    </react_native_1.View>
                                </react_native_1.View>
                            </react_native_1.TouchableOpacity>);
            }}/>
                </react_native_1.View>)}

            <react_native_paper_1.FAB style={styles_1.layout.fab} icon="plus" color={'white'} onPress={function () { return navigation.navigate(routes_1["default"].addRecord.name, { formId: formId, recordId: (0, uuid_1["default"])() }); }}/>
        </react_native_1.View>);
};
exports["default"] = RecordsScreen;
