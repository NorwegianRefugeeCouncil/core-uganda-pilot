"use strict";
var _a, _b;
exports.__esModule = true;
exports.host = void 0;
var expo_constants_1 = require("expo-constants");
exports.host = (_b = (_a = expo_constants_1["default"].manifest) === null || _a === void 0 ? void 0 : _a.extra) === null || _b === void 0 ? void 0 : _b.server_hostname;
