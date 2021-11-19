"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var search_white_png_1 = require("../../assets/png/search_white.png");
var symbol_filter_png_1 = require("../../assets/png/symbol_filter.png");
var symbol_individuals_png_1 = require("../../assets/png/symbol_individuals.png");
var styles_1 = require("../styles");
var routes_1 = require("../constants/routes");
var NavigationDropdown_1 = require("./NavigationDropdown");
var NavigationBar = function (_a) {
    var navigation = _a.navigation, back = _a.back, options = _a.options;
    var _b = react_1["default"].useState(false), visible = _b[0], setVisible = _b[1];
    var openMenu = function () { return setVisible(true); };
    var closeMenu = function () { return setVisible(false); };
    return (<react_native_paper_1.Appbar.Header>
      {back ? <react_native_paper_1.Appbar.BackAction onPress={navigation.goBack}/> : null}
      <react_native_paper_1.Appbar.Action icon={symbol_individuals_png_1["default"]} accessibilityLabel={routes_1["default"].addRecord.title} onPress={function () { return navigation.navigate(routes_1["default"].addRecord.name); }}/>
      <react_native_paper_1.Appbar.Action icon={search_white_png_1["default"]}/>
      <react_native_paper_1.Appbar.Content title={options.title} titleStyle={styles_1.common.textCentered} onPress={function () { return navigation.navigate(routes_1["default"].forms.name); }}/>
      <react_native_paper_1.Appbar.Action icon={symbol_filter_png_1["default"]}/>
      <NavigationDropdown_1["default"] visible={visible} closeMenu={closeMenu} openMenu={openMenu} navigation={navigation}/>
    </react_native_paper_1.Appbar.Header>);
};
exports["default"] = NavigationBar;
