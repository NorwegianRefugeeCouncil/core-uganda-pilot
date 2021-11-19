"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("../../constants/theme");
var TextInput = function (_a) {
    var fieldDefinition = _a.fieldDefinition, style = _a.style, value = _a.value, onChange = _a.onChange, onBlur = _a.onBlur, error = _a.error, invalid = _a.invalid, isTouched = _a.isTouched, isDirty = _a.isDirty, isMultiple = _a.isMultiple, isQuantity = _a.isQuantity;
    return (<react_native_1.View style={style}>
            {fieldDefinition.name && <react_native_paper_1.Text theme={theme_1.darkTheme}>{fieldDefinition.name}</react_native_paper_1.Text>}
            {fieldDefinition.description && (<react_native_paper_1.Text theme={theme_1.darkTheme} style={{ fontSize: 10 }}>
                    {fieldDefinition.description}
                </react_native_paper_1.Text>)}
            <react_native_paper_1.TextInput onChangeText={onChange} value={value} onBlur={onBlur} error={isTouched && isDirty && error} multiline={isMultiple} keyboardType={isQuantity ? "numeric" : "default"}/>
            {isTouched && isDirty && error && (<react_native_paper_1.Text>
                    {error.message == '' ? 'invalid' : error.message}
                </react_native_paper_1.Text>)}
        </react_native_1.View>);
};
exports["default"] = TextInput;
