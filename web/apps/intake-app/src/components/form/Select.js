"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("../../constants/theme");
var picker_1 = require("@react-native-picker/picker");
var Select = function (_a) {
    var _b;
    var fieldDefinition = _a.fieldDefinition, style = _a.style, value = _a.value, onChange = _a.onChange;
    var _c = react_1["default"].useState(value), selectedValue = _c[0], setSelectedValue = _c[1];
    return (<react_native_1.View style={style}>
            {fieldDefinition.name && <react_native_paper_1.Text theme={theme_1.darkTheme}>{fieldDefinition.name}</react_native_paper_1.Text>}
            {fieldDefinition.description && (<react_native_paper_1.Text theme={theme_1.darkTheme} style={{ fontSize: 10 }}>
                    {fieldDefinition.description}
                </react_native_paper_1.Text>)}
            <picker_1.Picker selectedValue={selectedValue} style={{ height: 50, width: 150 }} onValueChange={function (itemValue) {
            setSelectedValue(itemValue);
            onChange(itemValue);
        }}>
                {(_b = fieldDefinition.options) === null || _b === void 0 ? void 0 : _b.map(function (option) { return (<picker_1.Picker.Item key={option} label={option} value={option}/>); })}
            </picker_1.Picker>
        </react_native_1.View>);
};
exports["default"] = Select;
