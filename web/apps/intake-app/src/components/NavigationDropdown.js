"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var symbol_more_png_1 = require("../../assets/png/symbol_more.png");
var theme_1 = require("../constants/theme");
var routes_1 = require("../constants/routes");
var NavigationDropdown = function (_a) {
    var visible = _a.visible, closeMenu = _a.closeMenu, openMenu = _a.openMenu, navigation = _a.navigation;
    return (<react_native_paper_1.Menu visible={visible} onDismiss={closeMenu} anchor={<react_native_paper_1.Appbar.Action icon={symbol_more_png_1["default"]} onPress={openMenu} color={theme_1["default"].colors.white}/>}>
      <react_native_paper_1.Menu.Item title="Cases" onPress={function () { return navigation.navigate(routes_1["default"].cases.name); }}/>
      <react_native_paper_1.Menu.Item onPress={function () {
            console.log('Option 2 was pressed');
        }} title="Option 2"/>
      <react_native_paper_1.Menu.Item onPress={function () {
            console.log('Option 3 was pressed');
        }} title="Option 3" disabled/>
    </react_native_paper_1.Menu>);
};
exports["default"] = NavigationDropdown;
