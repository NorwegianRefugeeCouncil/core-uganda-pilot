"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("../../constants/theme");
var Switch = function (_a) {
    var fieldDefinition = _a.fieldDefinition, style = _a.style, value = _a.value, onChange = _a.onChange;
    return (<react_native_1.View style={style}>
            {fieldDefinition.label && (<react_native_paper_1.Text theme={theme_1.darkTheme}>
                    {fieldDefinition.label[0].value}
                </react_native_paper_1.Text>)}
            {fieldDefinition.description &&
            <react_native_paper_1.Text theme={theme_1.darkTheme} style={{ fontSize: 10 }}>
                {fieldDefinition.description[0].value}
            </react_native_paper_1.Text>}
            <react_native_paper_1.Switch value={value} onValueChange={onChange}/>
        </react_native_1.View>);
};
exports["default"] = Switch;
