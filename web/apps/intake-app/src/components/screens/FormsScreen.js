"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var styles_1 = require("../../styles");
var react_native_1 = require("react-native");
var clients_1 = require("../../utils/clients");
var routes_1 = require("../../constants/routes");
var FormsScreen = function (_a) {
    var navigation = _a.navigation;
    var _b = react_1["default"].useState(), forms = _b[0], setForms = _b[1];
    var _c = react_1["default"].useState(true), isLoading = _c[0], setIsLoading = _c[1];
    var client = (0, clients_1["default"])();
    react_1["default"].useEffect(function () {
        client.listForms({})
            .then(function (data) {
            var _a;
            setForms((_a = data.response) === null || _a === void 0 ? void 0 : _a.items);
            setIsLoading(false);
        });
    }, [client]);
    return (<react_native_1.View style={styles_1.layout.body}>
            <react_native_paper_1.Title>{routes_1["default"].forms.title}</react_native_paper_1.Title>
            {!isLoading && (<react_native_1.FlatList style={{ flex: 1, width: '100%' }} data={forms} renderItem={function (_a) {
                var item = _a.item, index = _a.index, separators = _a.separators;
                return (<react_native_1.TouchableOpacity key={index} onPress={function () { return navigation.navigate(routes_1["default"].records.name, {
                        formId: item.id,
                        databaseId: item.databaseId
                    }); }}>
                            <react_native_1.View style={{ flexDirection: 'row', flex: 1 }}>
                                <react_native_1.View style={{ justifyContent: 'center', paddingRight: 12 }}>
                                    <react_native_1.Text>{item.code}</react_native_1.Text>
                                </react_native_1.View>
                                <react_native_1.View style={{ justifyContent: 'center' }}>
                                    <react_native_1.Text>{item.name}</react_native_1.Text>
                                </react_native_1.View>
                            </react_native_1.View>
                        </react_native_1.TouchableOpacity>);
            }}/>)}
        </react_native_1.View>);
};
exports["default"] = FormsScreen;
