"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("../../constants/theme");
var lodash_1 = require("lodash");
// FIXME: reacts only on third click
var CheckBox = function (_a) {
    var fieldDefinition = _a.fieldDefinition, style = _a.style, onChange = _a.onChange;
    var _b = react_1["default"].useState(fieldDefinition.checkboxOptions.map(function (o) {
        if (fieldDefinition.value == null) {
            return o.value == 'yes';
        }
        else {
            return !!lodash_1["default"].find(fieldDefinition.value, function (e) {
                return e == o.value;
            });
        }
    })), isChecked = _b[0], setIsChecked = _b[1];
    return (<react_native_1.View style={style}>
            {fieldDefinition.label && (<react_native_1.Text>
                    {fieldDefinition.label[0].value}
                </react_native_1.Text>)}
            {fieldDefinition.checkboxOptions.map(function (option, i) { return (<react_native_1.View key={option.label[0].value}>
                    <react_native_1.View style={{
                display: 'flex',
                flexDirection: "row",
                alignItems: "center"
            }}>
                        <react_native_paper_1.Checkbox theme={theme_1.darkTheme} status={isChecked[i] ? 'checked' : 'unchecked'} onPress={function () {
                var isCheckedtmp = isChecked;
                isCheckedtmp[i] = !isCheckedtmp[i];
                setIsChecked(isCheckedtmp);
                onChange(isCheckedtmp[i]);
            }}/>
                        {option.label && (<react_native_1.Text>
                                {option.label[0].value}
                            </react_native_1.Text>)}
                    </react_native_1.View>
                </react_native_1.View>); })}
        </react_native_1.View>);
};
exports["default"] = CheckBox;
