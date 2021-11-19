"use strict";
exports.__esModule = true;
var axios_1 = require("axios");
var responses_1 = require("../utils/responses");
var client = /** @class */ (function () {
    function client(address, axiosInstance) {
        if (address === void 0) { address = 'https://localhost:9001/'; }
        if (axiosInstance === void 0) { axiosInstance = axios_1["default"].create(); }
        this.address = address;
        this.axiosInstance = axiosInstance;
        this.apiDomain = 'apis/admin.nrc.no';
    }
    client.prototype["do"] = function (request, url, method, data, expectStatusCode, options) {
        var headers = {
            "Accept": "application/json"
        };
        if (options === null || options === void 0 ? void 0 : options.headers) {
            headers = options === null || options === void 0 ? void 0 : options.headers;
        }
        return axios_1["default"].request({
            method: method,
            url: url,
            data: data,
            responseType: "json",
            headers: headers,
            withCredentials: true
        }).then(function (value) {
            return (0, responses_1.clientResponse)(value, request, expectStatusCode);
        })["catch"](function (err) {
            return {
                request: request,
                response: undefined,
                status: "500 Internal Server Error",
                statusCode: 500,
                error: err.message,
                success: false
            };
        });
    };
    client.prototype.createIdentityProvider = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/identityproviders", "post", request.object, 200);
    };
    client.prototype.createOrganization = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/organizations", "post", request.object, 200);
    };
    client.prototype.createOAuth2Client = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/clients", "post", request.object, 200);
    };
    client.prototype.deleteOAuth2Client = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/clients/" + request.id, "delete", undefined, 204);
    };
    client.prototype.getIdentityProvider = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/identityproviders/" + request.id, "get", undefined, 200);
    };
    client.prototype.getOrganization = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/organizations/" + request.id, "get", undefined, 200);
    };
    client.prototype.getOAuth2Client = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/clients/" + request.id, "get", undefined, 200);
    };
    client.prototype.listIdentityProviders = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/identityproviders?organizationId=" + request.organizationId, "get", undefined, 200);
    };
    client.prototype.listOAuth2Clients = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/clients", "get", undefined, 200);
    };
    client.prototype.listOrganizations = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/organizations", "get", undefined, 200);
    };
    client.prototype.updateIdentityProvider = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/identityproviders/" + request.object.id, "put", request.object, 200);
    };
    client.prototype.updateOAuth2Client = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/clients/" + request.object.id, "put", request.object, 200);
    };
    client.prototype.getSession = function (request) {
        return this["do"](request, this.address + "/" + this.apiDomain + "/v1/oidc/session", "get", undefined, 200, { headers: {} });
    };
    return client;
}());
exports["default"] = client;
