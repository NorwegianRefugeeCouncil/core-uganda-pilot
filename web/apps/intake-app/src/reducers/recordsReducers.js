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
exports.recordsReducer = exports.RECORD_ACTIONS = exports.initialRecordsState = void 0;
var lodash_1 = require("lodash");
exports.initialRecordsState = {
    formsById: {}
};
var RECORD_ACTIONS;
(function (RECORD_ACTIONS) {
    RECORD_ACTIONS["GET_RECORDS"] = "GET_RECORDS";
    RECORD_ACTIONS["GET_LOCAL_RECORDS"] = "GET_LOCAL_RECORDS";
    RECORD_ACTIONS["ADD_LOCAL_RECORD"] = "ADD_LOCAL_RECORD";
})(RECORD_ACTIONS = exports.RECORD_ACTIONS || (exports.RECORD_ACTIONS = {}));
var recordsReducer = function (state, action) {
    var _a, _b, _c;
    var formId = action.payload.formId;
    switch (action.type) {
        case RECORD_ACTIONS.GET_RECORDS:
            var records = action.payload.records;
            return {
                formsById: __assign(__assign({}, state.formsById), (_a = {}, _a[formId] = __assign(__assign({}, state.formsById[formId]), { records: records, recordsById: lodash_1["default"].keyBy(records, 'id') }), _a))
            };
        case RECORD_ACTIONS.GET_LOCAL_RECORDS:
            var localRecords = action.payload.localRecords;
            return {
                formsById: __assign(__assign({}, state.formsById), (_b = {}, _b[formId] = __assign(__assign({}, state.formsById[formId]), { localRecords: localRecords }), _b))
            };
        case RECORD_ACTIONS.ADD_LOCAL_RECORD:
            var currentList = state.formsById[formId].localRecords || [];
            var newList = currentList.concat([action.payload.localRecord]);
            return {
                formsById: __assign(__assign({}, state.formsById), (_c = {}, _c[formId] = __assign(__assign({}, state.formsById[formId]), { localRecords: newList }), _c))
            };
        default:
            throw new Error();
    }
};
exports.recordsReducer = recordsReducer;
