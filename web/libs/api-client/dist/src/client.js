"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.defaultClient = void 0;
const axios_1 = __importDefault(require("axios"));
const responses_1 = require("./utils/responses");
class client {
    constructor(address = 'https://core.dev:8443', axiosInstance = axios_1.default.create()) {
        this.address = address;
        this.axiosInstance = axiosInstance;
    }
    do(request, url, method, data, expectStatusCode) {
        let headers = {
            "Accept": "application/json",
        };
        return this.axiosInstance.request({
            responseType: "json",
            method,
            url,
            data,
            headers,
            withCredentials: true,
        }).then(value => {
            return (0, responses_1.clientResponse)(value, request, expectStatusCode);
        }).catch((err) => {
            var _a, _b;
            return {
                request: request,
                response: undefined,
                status: (_a = err.response) === null || _a === void 0 ? void 0 : _a.statusText,
                statusCode: (_b = err.response) === null || _b === void 0 ? void 0 : _b.status,
                error: err.response,
                success: false,
            };
        });
    }
    createDatabase(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/databases`, "post", request.object, 200);
    }
    createFolder(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/folders`, "post", request.object, 200);
    }
    createForm(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/forms`, "post", request.object, 200);
    }
    createRecord(request) {
        const url = `${this.address}/apis/core.nrc.no/v1/records`;
        return this.do(request, url, "post", request.object, 200);
    }
    listDatabases(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/databases`, "get", undefined, 200);
    }
    listFolders(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/folders`, "get", undefined, 200);
    }
    listForms(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/forms`, "get", undefined, 200);
    }
    listRecords(request) {
        const { databaseId, formId } = request;
        const url = `${this.address}/apis/core.nrc.no/v1/records?databaseId=${databaseId}&formId=${formId}`;
        return this.do(request, url, "get", undefined, 200);
    }
    getForm(request) {
        return this.do(request, `${this.address}/apis/core.nrc.no/v1/forms/${request.id}`, "get", undefined, 200);
    }
    getRecord(request) {
        const { databaseId, formId, recordId } = request;
        const url = `${this.address}/apis/core.nrc.no/v1/records/${recordId}?databaseId=${databaseId}&formId=${formId}`;
        return this.do(request, url, "get", undefined, 200);
    }
}
exports.default = client;
exports.defaultClient = new client();
