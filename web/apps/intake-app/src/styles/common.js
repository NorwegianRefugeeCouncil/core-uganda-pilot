"use strict";
exports.__esModule = true;
var react_native_1 = require("react-native");
var theme_1 = require("../constants/theme");
exports["default"] = react_native_1.StyleSheet.create({
    textCentered: {
        alignSelf: 'center'
    },
    darkBackground: {
        backgroundColor: theme_1["default"].colors.backdrop,
        color: theme_1["default"].colors.text
    }
});
