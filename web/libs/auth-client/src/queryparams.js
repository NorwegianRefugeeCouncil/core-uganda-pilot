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
exports.getQueryParams = exports.buildQueryString = void 0;
var qs_1 = require("qs");
function buildQueryString(input) {
    return qs_1["default"].stringify(input);
}
exports.buildQueryString = buildQueryString;
function getQueryParams(url) {
    var _a;
    var parts = url.split('#');
    var hash = parts[1];
    var partsWithoutHash = parts[0].split('?');
    var queryString = partsWithoutHash[partsWithoutHash.length - 1];
    // Get query string (?hello=world)
    var parsedSearch = qs_1["default"].parse(queryString, { parseArrays: false });
    // Pull errorCode off of params
    var errorCode = ((_a = parsedSearch.errorCode) !== null && _a !== void 0 ? _a : null);
    delete parsedSearch.errorCode;
    // Get hash (#abc=example)
    var parsedHash = {};
    if (parts[1]) {
        parsedHash = qs_1["default"].parse(hash);
    }
    // Merge search and hash
    var params = __assign(__assign({}, parsedSearch), parsedHash);
    return {
        errorCode: errorCode,
        params: params
    };
}
exports.getQueryParams = getQueryParams;
