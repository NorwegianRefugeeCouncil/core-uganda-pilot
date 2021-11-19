"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getQueryParams = exports.buildQueryString = void 0;
const qs_1 = __importDefault(require("qs"));
function buildQueryString(input) {
    return qs_1.default.stringify(input);
}
exports.buildQueryString = buildQueryString;
function getQueryParams(url) {
    var _a;
    const parts = url.split('#');
    const hash = parts[1];
    const partsWithoutHash = parts[0].split('?');
    const queryString = partsWithoutHash[partsWithoutHash.length - 1];
    // Get query string (?hello=world)
    const parsedSearch = qs_1.default.parse(queryString, { parseArrays: false });
    // Pull errorCode off of params
    const errorCode = ((_a = parsedSearch.errorCode) !== null && _a !== void 0 ? _a : null);
    delete parsedSearch.errorCode;
    // Get hash (#abc=example)
    let parsedHash = {};
    if (parts[1]) {
        parsedHash = qs_1.default.parse(hash);
    }
    // Merge search and hash
    const params = Object.assign(Object.assign({}, parsedSearch), parsedHash);
    return {
        errorCode,
        params,
    };
}
exports.getQueryParams = getQueryParams;
