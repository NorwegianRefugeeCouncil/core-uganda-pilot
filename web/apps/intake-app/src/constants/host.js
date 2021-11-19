"use strict";
var _a, _b, _c, _d;
exports.__esModule = true;
var react_native_1 = require("react-native");
var expo_constants_1 = require("expo-constants");
var host = react_native_1.Platform.OS === 'web' ?
    (_b = (_a = expo_constants_1["default"].manifest) === null || _a === void 0 ? void 0 : _a.extra) === null || _b === void 0 ? void 0 : _b.server_default_hostname :
    (_d = (_c = expo_constants_1["default"].manifest) === null || _c === void 0 ? void 0 : _c.extra) === null || _d === void 0 ? void 0 : _d.server_hostname;
exports["default"] = host;
