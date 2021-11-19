"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("../../constants/theme");
var RadioButtons = function (_a) {
    var fieldDefinition = _a.fieldDefinition;
    var _b = react_1["default"].useState(fieldDefinition.value[0]), checked = _b[0], setChecked = _b[1];
    return (<react_native_1.View>
            {fieldDefinition.label && <react_native_paper_1.Text theme={theme_1.darkTheme}>{fieldDefinition.label[0].value}</react_native_paper_1.Text>}
            {fieldDefinition.description &&
            <react_native_paper_1.Text theme={theme_1.darkTheme} style={{ fontSize: 10 }}>{fieldDefinition.description[0].value}</react_native_paper_1.Text>}

            {fieldDefinition.options.map(function (option) { return (<react_native_paper_1.RadioButton key={option[0].value} value={option[0].value} status={checked === 'first' ? 'checked' : 'unchecked'} onPress={function () { return setChecked('first'); }}/>); })}
        </react_native_1.View>);
};
exports["default"] = RadioButtons;
