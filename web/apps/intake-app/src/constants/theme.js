"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
exports.__esModule = true;
exports.NavigationTheme = exports.darkTheme = void 0;
var react_native_paper_1 = require("react-native-paper");
var native_1 = require("@react-navigation/native");
var lightTheme = __assign(__assign({}, react_native_paper_1.DefaultTheme), { colors: {
        primary: '#24303E',
        accent: '#FF7602',
        background: '#E5E5E5',
        surface: '#E5E5E5',
        error: '#CE3616',
        text: '#000000',
        onSurface: '#FF7602',
        disabled: '#808080',
        placeholder: '#808080',
        backdrop: '#24303E',
        notification: '#FFF0E6',
        nrc: '#FF7602',
        green: '#47914A',
        blue: '#00ADD0',
        yellow: '#FDC82F',
        red: '#CE3616',
        white: '#FFFFFF'
    } });
exports["default"] = lightTheme;
exports.darkTheme = __assign(__assign(__assign({}, react_native_paper_1.DefaultTheme), lightTheme), { dark: true, colors: __assign(__assign({}, lightTheme.colors), { primary: '#24303E', accent: '#FF7602', background: '#24303E', surface: '#808080', text: '#FFFFFF', onSurface: '#FF7602', backdrop: '#808080', notification: '#808080' }) });
exports.NavigationTheme = __assign(__assign({}, native_1.DefaultTheme), { dark: false, colors: __assign(__assign(__assign({}, native_1.DefaultTheme.colors), lightTheme.colors), { card: '#24303E', text: '#FFFFFF', border: '#24303E' }) });
