"use strict";
exports.__esModule = true;
exports.defaultClient = void 0;
var axios_1 = require("axios");
var responses_1 = require("./utils/responses");
var client = /** @class */ (function () {
    function client(address, axiosInstance) {
        if (address === void 0) { address = 'https://core.dev:8443'; }
        if (axiosInstance === void 0) { axiosInstance = axios_1["default"].create(); }
        this.address = address;
        this.axiosInstance = axiosInstance;
    }
    client.prototype["do"] = function (request, url, method, data, expectStatusCode) {
        var headers = {
            "Accept": "application/json"
        };
        return this.axiosInstance.request({
            responseType: "json",
            method: method,
            url: url,
            data: data,
            headers: headers,
            withCredentials: true
        }).then(function (value) {
            return (0, responses_1.clientResponse)(value, request, expectStatusCode);
        })["catch"](function (err) {
            var _a, _b;
            return {
                request: request,
                response: undefined,
                status: (_a = err.response) === null || _a === void 0 ? void 0 : _a.statusText,
                statusCode: (_b = err.response) === null || _b === void 0 ? void 0 : _b.status,
                error: err.response,
                success: false
            };
        });
    };
    client.prototype.createDatabase = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/databases", "post", request.object, 200);
    };
    client.prototype.createFolder = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/folders", "post", request.object, 200);
    };
    client.prototype.createForm = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/forms", "post", request.object, 200);
    };
    client.prototype.createRecord = function (request) {
        var url = this.address + "/apis/core.nrc.no/v1/records";
        return this["do"](request, url, "post", request.object, 200);
    };
    client.prototype.listDatabases = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/databases", "get", undefined, 200);
    };
    client.prototype.listFolders = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/folders", "get", undefined, 200);
    };
    client.prototype.listForms = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/forms", "get", undefined, 200);
    };
    client.prototype.listRecords = function (request) {
        var databaseId = request.databaseId, formId = request.formId;
        var url = this.address + "/apis/core.nrc.no/v1/records?databaseId=" + databaseId + "&formId=" + formId;
        return this["do"](request, url, "get", undefined, 200);
    };
    client.prototype.getForm = function (request) {
        return this["do"](request, this.address + "/apis/core.nrc.no/v1/forms/" + request.id, "get", undefined, 200);
    };
    client.prototype.getRecord = function (request) {
        var databaseId = request.databaseId, formId = request.formId, recordId = request.recordId;
        var url = this.address + "/apis/core.nrc.no/v1/records/" + recordId + "?databaseId=" + databaseId + "&formId=" + formId;
        return this["do"](request, url, "get", undefined, 200);
    };
    return client;
}());
exports["default"] = client;
exports.defaultClient = new client();
